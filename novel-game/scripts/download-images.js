const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');

// Create directories if they don't exist
const createDirectories = () => {
  const dirs = [
    path.join(__dirname, '../public/assets/characters'),
    path.join(__dirname, '../public/assets/backgrounds'),
  ];

  dirs.forEach(dir => {
    if (!fs.existsSync(dir)) {
      fs.mkdirSync(dir, { recursive: true });
      console.log(`Created directory: ${dir}`);
    }
  });
};

// Download an image using curl (more reliable for redirects)
const downloadImage = (url, filepath) => {
  return new Promise((resolve, reject) => {
    try {
      // Use curl to download the image with follow redirects
      execSync(`curl -L "${url}" -o "${filepath}"`);
      console.log(`Downloaded: ${filepath}`);
      resolve();
    } catch (error) {
      reject(`Error downloading ${url}: ${error.message}`);
    }
  });
};

// Character images to download
const characterImages = [
  { name: 'tanaka_hiroshi', width: 800, height: 600, id: 1 },
  { name: 'kim_min_ji', width: 800, height: 600, id: 2 },
  { name: 'garcia_carlos', width: 800, height: 600, id: 3 },
  { name: 'yamamoto_sensei', width: 800, height: 600, id: 4 },
  { name: 'nakamura_yuki', width: 800, height: 600, id: 5 },
  { name: 'suzuki_kenji', width: 800, height: 600, id: 6 },
  { name: 'watanabe_akiko', width: 800, height: 600, id: 7 },
  { name: 'alex_thompson', width: 800, height: 600, id: 8 },
];

// Background images to download
const backgroundImages = [
  { name: 'post_office', width: 1920, height: 1080, id: 10 },
  { name: 'cafe', width: 1920, height: 1080, id: 11 },
  { name: 'classroom', width: 1920, height: 1080, id: 12 },
  { name: 'apartment', width: 1920, height: 1080, id: 13 },
  { name: 'corner_store', width: 1920, height: 1080, id: 14 },
];

// Main function to download all images
const downloadAllImages = async () => {
  try {
    createDirectories();

    // Download character images
    for (const img of characterImages) {
      const url = `https://picsum.photos/id/${img.id}/${img.width}/${img.height}`;
      const filepath = path.join(__dirname, `../public/assets/characters/${img.name}.jpg`);
      await downloadImage(url, filepath);
    }

    // Download background images
    for (const img of backgroundImages) {
      const url = `https://picsum.photos/id/${img.id}/${img.width}/${img.height}`;
      const filepath = path.join(__dirname, `../public/assets/backgrounds/${img.name}.jpg`);
      await downloadImage(url, filepath);
    }

    console.log('All images downloaded successfully!');
  } catch (error) {
    console.error('Error downloading images:', error);
  }
};

// Run the script
downloadAllImages(); 