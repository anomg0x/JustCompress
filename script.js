document.addEventListener('DOMContentLoaded', () => {
    const dropZone = document.getElementById('drop-zone');
    const fileInput = document.getElementById('file-input');
    const controlsSection = document.getElementById('controls-section');
    const previewSection = document.getElementById('preview-section');
    
    const qualitySlider = document.getElementById('quality-slider');
    const qualityValue = document.getElementById('quality-value');
    const sliderContainer = document.getElementById('slider-container');
    const losslessToggle = document.getElementById('lossless-toggle');
    const modeText = document.getElementById('mode-text');
    const serverToggle = document.getElementById('server-toggle');
    const engineText = document.getElementById('engine-text');
    
    const originalImg = document.getElementById('original-img');
    const compressedImg = document.getElementById('compressed-img');
    const originalSizeTxt = document.getElementById('original-size');
    const compressedSizeTxt = document.getElementById('compressed-size');
    const savingsBadge = document.getElementById('savings-badge');
    const downloadBtn = document.getElementById('download-btn');

    let currentFile = null;
    let originalImageObj = new Image();

    const GO_SERVER_URL = '';

    // Engine Toggle Event
    serverToggle.addEventListener('change', (e) => {
        if (e.target.checked) {
            engineText.textContent = "Go Server (Backend)";
            engineText.classList.add('text-purple-400');
        } else {
            engineText.textContent = "Client-side (Browser)";
            engineText.classList.remove('text-purple-400');
        }
        if (currentFile) compressImage();
    });

    // Drag and Drop Events
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, preventDefaults, false);
    });

    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    ['dragenter', 'dragover'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => dropZone.classList.add('drag-over'), false);
    });

    ['dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => dropZone.classList.remove('drag-over'), false);
    });

    dropZone.addEventListener('drop', (e) => {
        handleFiles(e.dataTransfer.files);
    });

    // Click to upload
    dropZone.addEventListener('click', () => fileInput.click());
    fileInput.addEventListener('change', function() { handleFiles(this.files); });

    // Slider Event
    qualitySlider.addEventListener('input', (e) => {
        qualityValue.textContent = `${e.target.value}%`;
        if (currentFile) compressImage();
    });

    // Lossless Toggle Event
    losslessToggle.addEventListener('change', (e) => {
        if (e.target.checked) {
            modeText.textContent = "Lossless (100% Quality)";
            sliderContainer.classList.add('opacity-50', 'pointer-events-none');
            qualityValue.textContent = "100%";
        } else {
            modeText.textContent = "Lossy (Smaller size)";
            sliderContainer.classList.remove('opacity-50', 'pointer-events-none');
            qualityValue.textContent = `${qualitySlider.value}%`;
        }
        if (currentFile) compressImage();
    });

    function handleFiles(files) {
        if (files.length === 0) return;
        
        const file = files[0];
        if (!file.type.match(/image\/(jpeg|png|webp)/)) {
            alert('Please upload a valid image file (JPEG, PNG, or WEBP).');
            return;
        }

        currentFile = file;
        const reader = new FileReader();

        reader.onload = (e) => {
            originalImageObj.src = e.target.result;
            originalImageObj.onload = () => {
                // UI Updates
                originalImg.src = originalImageObj.src;
                originalSizeTxt.textContent = formatBytes(currentFile.size);
                
                controlsSection.classList.remove('hidden');
                previewSection.classList.remove('hidden');
                
                compressImage(); // Initial compression
            };
        };
        reader.readAsDataURL(file);
    }

    function compressImage() {
        if (serverToggle.checked) {
            compressWithServer();
        } else {
            compressClientSide();
        }
    }

    async function compressWithServer() {
        const isLossless = losslessToggle.checked;
        const quality = isLossless ? 100 : parseInt(qualitySlider.value);

        // Get format from file type
        let format = 'jpeg';
        if (currentFile.type === 'image/png') format = 'png';
        else if (currentFile.type === 'image/webp') format = 'jpeg'; // Go can't encode webp

        const formData = new FormData();
        formData.append('image', currentFile);
        formData.append('quality', quality.toString());
        formData.append('format', format);
        formData.append('lossless', isLossless.toString());

        // Show loading state
        compressedSizeTxt.textContent = 'Compressing...';
        savingsBadge.textContent = '⏳';
        savingsBadge.className = "bg-yellow-900 text-yellow-300 px-2 py-1 rounded text-xs font-bold";

        try {
            const response = await fetch(`${GO_SERVER_URL}/api/compress`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                const err = await response.json();
                throw new Error(err.error || 'Server error');
            }

            const blob = await response.blob();
            const compressedUrl = URL.createObjectURL(blob);
            const compressedSize = blob.size;

            compressedImg.src = compressedUrl;
            compressedSizeTxt.textContent = formatBytes(compressedSize);
            updateSavingsBadge(currentFile.size, compressedSize);

            // Setup download
            downloadBtn.onclick = () => {
                const link = document.createElement('a');
                const ext = format === 'png' ? 'png' : 'jpg';
                const suffix = isLossless ? '_lossless' : '_compressed';
                link.download = `JustCompress_${currentFile.name.replace(/\.[^/.]+$/, "")}${suffix}_go.${ext}`;
                link.href = compressedUrl;
                link.click();
            };

            // Server response headers with metadata
            const duration = response.headers.get('X-Duration');
            if (duration) {
                savingsBadge.title = `Server compression time: ${duration}`;
            }

        } catch (err) {
            compressedSizeTxt.textContent = 'Error';
            savingsBadge.textContent = '❌';
            savingsBadge.className = "bg-red-900 text-red-300 px-2 py-1 rounded text-xs font-bold";
            console.error('Go server compression error:', err);

            // Fallback: inform the user
            alert(`Go server error: ${err.message}\n\nMake sure the Go server is running:\ncd server && ./justcompress-server.exe`);
        }
    }

    function compressClientSide() {
        // Если включен Lossless, ставим максимальное качество
        const isLossless = losslessToggle.checked;
        const quality = isLossless ? 1.0 : parseInt(qualitySlider.value) / 100;
        
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');

        canvas.width = originalImageObj.width;
        canvas.height = originalImageObj.height;
        ctx.drawImage(originalImageObj, 0, 0, canvas.width, canvas.height);

        // Используем оригинальный тип файла
        let mimeType = currentFile.type;
        
        // Генерация обработанной картинки
        const compressedDataUrl = canvas.toDataURL(mimeType, quality);
        compressedImg.src = compressedDataUrl;

        // Расчет веса
        const base64Length = compressedDataUrl.length - (compressedDataUrl.indexOf(',') + 1);
        const padding = (compressedDataUrl.charAt(compressedDataUrl.length - 2) === '=') ? 2 : ((compressedDataUrl.charAt(compressedDataUrl.length - 1) === '=') ? 1 : 0);
        const compressedSize = (base64Length * 0.75) - padding;

        compressedSizeTxt.textContent = formatBytes(compressedSize);
        updateSavingsBadge(currentFile.size, compressedSize);

        // Настройка кнопки скачивания
        downloadBtn.onclick = () => {
            const link = document.createElement('a');
            const ext = mimeType === 'image/png' ? 'png' : (mimeType === 'image/webp' ? 'webp' : 'jpg');
            const suffix = isLossless ? '_lossless' : '_compressed';
            link.download = `JustCompress_${currentFile.name.replace(/\.[^/.]+$/, "")}${suffix}.${ext}`;
            link.href = compressedDataUrl;
            link.click();
        };

        // Отключаем UI ползунка для PNG в Lossy режиме (Canvas не сжимает PNG с потерями)
        if (mimeType === 'image/png' && !isLossless) {
             qualityValue.textContent = "Lossless format";
             sliderContainer.classList.add('opacity-50', 'pointer-events-none');
        } else if (mimeType !== 'image/png' && !isLossless) {
             sliderContainer.classList.remove('opacity-50', 'pointer-events-none');
             qualityValue.textContent = `${qualitySlider.value}%`;
        }
    }

    function updateSavingsBadge(original, compressed) {
        if (compressed >= original) {
            const increase = Math.round(((compressed / original) - 1) * 100);
            savingsBadge.textContent = `+${increase}%`;
            savingsBadge.className = "bg-red-900 text-red-300 px-2 py-1 rounded text-xs font-bold";
        } else {
            const savings = Math.round((1 - (compressed / original)) * 100);
            savingsBadge.textContent = `-${savings}%`;
            savingsBadge.className = "bg-green-900 text-green-300 px-2 py-1 rounded text-xs font-bold";
        }
    }

    function formatBytes(bytes, decimals = 2) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const dm = decimals < 0 ? 0 : decimals;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
    }
});
