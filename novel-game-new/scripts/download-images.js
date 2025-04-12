const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

// Configuration
const config = {
  characters: [
    { name: 'sensei', width: 500, height: 800, id: 1001 },
    { name: 'student', width: 500, height: 800, id: 1002 }
  ],
  backgrounds: [
    { name: 'classroom', width: 1920, height: 1080, id: 1003 },
    { name: 'school', width: 1920, height: 1080, id: 1004 },
    { name: 'library', width: 1920, height: 1080, id: 1005 }
  ]
};

// Ensure directories exist
const createDirIfNotExists = (dirPath) => {
  const fullPath = path.join(process.cwd(), '..', 'public', 'images', dirPath);
  if (!fs.existsSync(fullPath)) {
    fs.mkdirSync(fullPath, { recursive: true });
  }
};

// Download image using curl
const downloadImage = (type, item) => {
  const outputDir = path.join(process.cwd(), '..', 'public', 'images', type);
  const outputPath = path.join(outputDir, `${item.name}.jpg`);
  const url = `https://picsum.photos/id/${item.id}/${item.width}/${item.height}`;
  
  console.log(`Downloading ${type}/${item.name}.jpg from ${url}`);
  try {
    execSync(`curl -L "${url}" -o "${outputPath}"`);
    console.log(`Successfully downloaded ${type}/${item.name}.jpg`);
  } catch (error) {
    console.error(`Failed to download ${type}/${item.name}.jpg:`, error.message);
  }
};

// Main function
const downloadAllImages = async () => {
  // Create directories
  createDirIfNotExists('characters');
  createDirIfNotExists('backgrounds');

  // Download character images
  console.log('\nDownloading character images...');
  config.characters.forEach(char => downloadImage('characters', char));

  // Download background images
  console.log('\nDownloading background images...');
  config.backgrounds.forEach(bg => downloadImage('backgrounds', bg));

  console.log('\nDownload process completed!');
};

// Run the script
downloadAllImages().catch(console.error); 