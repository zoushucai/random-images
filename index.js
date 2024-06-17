const express = require('express');
const fs = require('fs');
const path = require('path');
const sharp = require('sharp');
const cors = require('cors');

const app = express();
const port = 2113;
const imagesDir = path.join(__dirname, 'images');
const imagesData = JSON.parse(fs.readFileSync('images_info.json', 'utf8'));

app.use(cors());
app.use('/images', express.static(imagesDir));

app.get('/random', async (req, res) => {
    try {
        let { sub, width, type, contains, index, device, json } = req.query;

        // Default values
        width = parseInt(width) || 1920;
        type = type || 'webp';
        contains = contains || '';
        device = device || 'pc';
        json = parseInt(json) || 0;

        // Filter images based on query params
        let filteredImages = imagesData.filter(image => (
            (!sub || image.sub === sub) && (!contains || image.file.includes(contains))
        ));

        if (filteredImages.length === 0) {
            return res.status(404).send('No images found.');
        }

        // Select random image index
        index = parseInt(index) || Math.floor(Math.random() * filteredImages.length);
        if (index < 0 || index >= filteredImages.length) {
            return res.status(400).send('Index out of range.');
        }

        const selectedImage = filteredImages[index];
        const imagePath = path.join(imagesDir, selectedImage.sub, selectedImage.file);

        let image = sharp(imagePath);
        const metadata = await image.metadata();

        if (type !== 'webp') {
            type = metadata.format; // Use original format if not webp
        }

        image = image.resize({ width });

        // Handle rotation for mobile devices
        if (['mobile', 'tablet', 'phone'].includes(device)) {
            // image = image.rotate(width < height ? 90 : 0);
            const metadata2 = await image.metadata();
            const newWidth = metadata2.height;
            const newHeight = metadata2.width;
            image = image.resize(newWidth, newHeight);
        }

        const buffer = await image.toBuffer();

        if (json === 1) {
            const imageData = {
                width: metadata.width,
                height: metadata.height,
                type,
                sub: selectedImage.sub,
                imageurl: selectedImage.file
            };
            res.json(imageData);
        } else {
            res.type(type);
            res.send(buffer);
        }
    } catch (error) {
        console.error('Error processing image:', error);
        res.status(500).send('Error processing image.');
    }
});

app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});