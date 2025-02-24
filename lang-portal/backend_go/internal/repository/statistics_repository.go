package repository

import (
	"database/sql"
	"time"

	"lang-portal/backend_go/internal/models"
)

type StatisticsRepository struct {
	db *sql.DB
}

func NewStatisticsRepository(db *sql.DB) *StatisticsRepository {
	return &StatisticsRepository{db: db}
}

func (r *StatisticsRepository) GetUserStats() (*models.UserStats, error) {
	var stats models.UserStats
	var lastStudy sql.NullTime

	err := r.db.QueryRow(`
		WITH study_stats AS (
			SELECT 
				SUM(CASE WHEN completed_at IS NOT NULL THEN 1 ELSE 0 END) as completed,
				COUNT(*) as total,
				AVG(CASE 
					WHEN accuracy_score IS NOT NULL THEN accuracy_score 
					WHEN confidence_score IS NOT NULL THEN confidence_score
					ELSE 0 END) as avg_score,
				MAX(created_at) as last_study
			FROM study_activities
		),
		time_stats AS (
			SELECT SUM(time_taken_seconds) / 60 as total_time
			FROM study_sessions
		),
		word_stats AS (
			SELECT COUNT(DISTINCT word_id) as words_learned
			FROM word_review_items
			WHERE correct = 1
		)
		SELECT 
			COALESCE(t.total_time, 0),
			s.total,
			s.completed,
			COALESCE(s.avg_score, 0),
			COALESCE(w.words_learned, 0),
			s.last_study
		FROM study_stats s
		LEFT JOIN time_stats t ON 1=1
		LEFT JOIN word_stats w ON 1=1
	`).Scan(
		&stats.TotalStudyTime,
		&stats.TotalActivities,
		&stats.CompletedActivities,
		&stats.AverageAccuracy,
		&stats.WordsLearned,
		&lastStudy,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if lastStudy.Valid {
		stats.LastStudyDate = &lastStudy.Time
	}

	return &stats, nil
}

func (r *StatisticsRepository) GetActivityStats() (map[string]models.ActivityStats, error) {
	rows, err := r.db.Query(`
		SELECT 
			type as activity_type,
			COUNT(*) as total,
			COUNT(CASE WHEN completed_at IS NOT NULL THEN 1 END) as completed,
			AVG(CASE 
				WHEN accuracy_score IS NOT NULL THEN accuracy_score 
				WHEN confidence_score IS NOT NULL THEN confidence_score
				ELSE 0 END) as avg_score,
			COALESCE(SUM(
				SELECT time_taken_seconds 
				FROM study_sessions 
				WHERE study_activity_id = a.id
			) / 60, 0) as total_time,
			MAX(created_at) as last_attempt
		FROM study_activities a
		GROUP BY type
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]models.ActivityStats)
	for rows.Next() {
		var (
			actType     string
			total       int
			completed   int
			avgScore    float64
			totalTime   int
			lastAttempt time.Time
		)

		err := rows.Scan(&actType, &total, &completed, &avgScore, &totalTime, &lastAttempt)
		if err != nil {
			return nil, err
		}

		stats[actType] = models.ActivityStats{
			ActivityType:    actType,
			CompletionRate:  float64(completed) / float64(total) * 100,
			AverageAccuracy: avgScore,
			TotalTime:       totalTime,
			LastAttempt:     lastAttempt,
		}
	}

	return stats, nil
}

func (r *StatisticsRepository) GetStudyProgress() (*models.StudyProgress, error) {
	// Get daily streak
	var streak int
	err := r.db.QueryRow(`
		WITH RECURSIVE dates AS (
			SELECT date(MAX(created_at)) as date
			FROM study_activities
			WHERE completed_at IS NOT NULL
			UNION ALL
			SELECT date(date, '-1 day')
			FROM dates
			WHERE EXISTS (
				SELECT 1 FROM study_activities 
				WHERE date(created_at) = date(dates.date, '-1 day')
				AND completed_at IS NOT NULL
			)
		)
		SELECT COUNT(*) FROM dates
	`).Scan(&streak)
	if err != nil {
		return nil, err
	}

	// Get weekly progress
	rows, err := r.db.Query(`
		SELECT 
			date(created_at) as study_date,
			SUM(time_taken_seconds) / 60 as minutes,
			COUNT(DISTINCT study_activity_id) as activities
		FROM study_sessions
		WHERE created_at >= date('now', '-7 days')
		GROUP BY date(created_at)
		ORDER BY study_date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progress []models.DailyProgress
	for rows.Next() {
		var dp models.DailyProgress
		err := rows.Scan(&dp.Date, &dp.Minutes, &dp.Activities)
		if err != nil {
			return nil, err
		}
		progress = append(progress, dp)
	}

	// Get activity type stats
	activityStats, err := r.GetActivityStats()
	if err != nil {
		return nil, err
	}

	return &models.StudyProgress{
		DailyStreak:    streak,
		WeeklyProgress: progress,
		ByActivityType: activityStats,
	}, nil
} 