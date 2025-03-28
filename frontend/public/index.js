"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
function uploadFile() {
    return __awaiter(this, void 0, void 0, function* () {
        const fileInput = document.getElementById('fileInput');
        const statusDiv = document.getElementById('status');
        if (!fileInput.files || fileInput.files.length === 0) {
            statusDiv.textContent = "Выберите файл!";
            return;
        }
        const formData = new FormData();
        formData.append('receipt', fileInput.files[0]);
        try {
            statusDiv.textContent = "Отправка файла...";
            const response = yield fetch('http://localhost:8080/api/upload', {
                method: 'POST',
                body: formData
            });
            if (!response.ok)
                throw new Error('Ошибка сервера');
            const result = yield response.json();
            statusDiv.textContent = `Успех: ${result.message}`;
        }
        catch (error) {
            if (error instanceof Error) {
                statusDiv.textContent = `Ошибка: ${error.message}`;
            }
            else {
                statusDiv.textContent = 'Неизвестная ошибка';
            }
        }
    });
}
