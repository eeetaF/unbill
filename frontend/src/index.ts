async function uploadFile(): Promise<void> {
    const fileInput = document.getElementById('fileInput') as HTMLInputElement;
    const statusDiv = document.getElementById('status') as HTMLElement;
    
    if (!fileInput.files || fileInput.files.length === 0) {
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
        if (error instanceof Error) {
            statusDiv.textContent = `Ошибка: ${error.message}`;
        } else {
            statusDiv.textContent = 'Неизвестная ошибка';
        }
    }
}
