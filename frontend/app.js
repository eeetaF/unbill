async function uploadFile() {
    const fileInput = document.getElementById('fileInput');
    const statusDiv = document.getElementById('status');
    
    if (!fileInput.files[0]) {
        statusDiv.textContent = "Выберите файл!";
        return;
    }

    const formData = new FormData();
    formData.append('receipt', fileInput.files[0]);

    try {
        statusDiv.textContent = "Отправка файла...";
        
        const response = await fetch('http://localhost:8080/api/upload', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) throw new Error('Ошибка сервера');
        
        const result = await response.json();
        statusDiv.textContent = `Успех: ${result.message}`;
        
    } catch (error) {
        statusDiv.textContent = `Ошибка: ${error.message}`;
    }
}