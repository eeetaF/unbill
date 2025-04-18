import base64
import logging
from pathlib import Path
from telegram import Update, InlineKeyboardButton, InlineKeyboardMarkup
from telegram.ext import Application, CommandHandler, MessageHandler, filters, ContextTypes, CallbackQueryHandler
import httpx
from dotenv import load_dotenv
import os 


load_dotenv()

BOT_TOKEN = os.getenv("TELEGRAM_TOKEN")
BACKEND_URL = "http://localhost:8080/api/upload_and_analyze"

logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    level=logging.INFO
)
logger = logging.getLogger(__name__)

async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    welcome_message = (
        "👋 Привет! Я — Unbill, бот, который помогает быстро и честно разделить чек между друзьями.\n\n"
        "Просто отправь фото чека 📷 — я распознаю все позиции и предложу, кто сколько должен.\n"
        "Никаких калькуляторов и споров — только удобство!\n\n"
        # "🥗 Обедали вместе?\n"
        # "🛍️ Покупали продукты на всех?\n"
        # "🍕 Заказали пиццу в складчину?\n\n"
        "Отправляй чек — и я всё рассчитаю ✨"
    )

    await update.message.reply_text(
        welcome_message
    )

async def handle_photo(update: Update, context: ContextTypes.DEFAULT_TYPE):
    try:
        photo_file = await update.message.photo[-1].get_file()
        file_ext = Path(photo_file.file_path).suffix[1:].lower()

        photo_extensions = [
            "jpg", "jpeg",  
            "png",    
            "bmp"      
            # "heic"         
            # ".webp",    
            # ".tiff", ".tif", 
            # ".gif"          
        ]
        
        if file_ext not in photo_extensions:
            print(file_ext)
            unsupported_format_message = (
                "😕 К сожалению, этот формат файла не поддерживается.\n"
                "Пожалуйста, отправь фото чека в одном из следующих форматов: JPG, PNG, HEIC."
            )
            await update.message.reply_text(unsupported_format_message)
            return

        async with httpx.AsyncClient() as client:
            
            image_data = await client.get(photo_file.file_path)
            image_data.raise_for_status()
            
            response = await client.post(
                BACKEND_URL,
                json={
                    "ext": file_ext,
                    "content": base64.b64encode(image_data.content).decode('utf-8')
                },
                timeout=60.0
            )
            
            if response.status_code == 200:
                try:
                    receipt_data = response.json()

                    context.user_data["filename"] = receipt_data["filename"]

                    message_text = format_receipt_message(receipt_data)
                    await update.message.reply_text(message_text, parse_mode="Markdown")

                    keyboard = [
                        [
                            InlineKeyboardButton("Поделить поровну", callback_data="split_even"),
                            InlineKeyboardButton("Выбрать вручную", callback_data="split_manual"),
                        ]
                    ]
                    reply_markup = InlineKeyboardMarkup(keyboard)

                    await update.message.reply_text(
                        "Как будем делить чек? 🤔",
                        reply_markup=reply_markup
                    )
                    
                except ValueError as e:
                    logger.error(f"Invalid JSON response: {e}")
                    await update.message.reply_text("❌ Ошибка формата ответа от сервера")
            else:
                await update.message.reply_text(f"❌ Ошибка сервера: {response.status_code}")


    except Exception as e:
        logger.error(f"Error: {str(e)}", exc_info=True)
        await update.message.reply_text("⚠️ Ошибка при обработке изображения")


async def handle_split_choice(update: Update, context: ContextTypes.DEFAULT_TYPE):
    query = update.callback_query
    await query.answer()

    if query.data == "split_even":
        await query.edit_message_text("🔢 На сколько человек делим счет?")
        context.user_data["awaiting_participant_count"] = True

    elif query.data == "split_manual":
        await query.edit_message_text("📝 Окей! Перейдём к ручному распределению.")
        # Реализация ручного режима позже

async def handle_participant_count(update: Update, context: ContextTypes.DEFAULT_TYPE):
    if not context.user_data.get("awaiting_participant_count"):
        return  # Игнорируем, если бот этого не ждал

    try:
        num_people = int(update.message.text.strip())

        if num_people <= 0:
            await update.message.reply_text("❌ Количество участников должно быть больше нуля. Попробуй ещё раз.")
            return

        filename = context.user_data.get("filename")
        if not filename:
            await update.message.reply_text("⚠️ Не удалось найти данные чека. Попробуй заново загрузить чек.")
            return

        # Убираем флаг ожидания, чтобы не перехватывать случайные сообщения
        context.user_data["awaiting_participant_count"] = False

        # Делаем запрос
        url = f"http://localhost:8080/api/split_equally/{filename}/{num_people}"
        print(url)
        async with httpx.AsyncClient() as client:
            response = await client.get(url, timeout=30.0)

        if response.status_code == 200:
            result = response.json()
            text = format_split_result(result)
            await update.message.reply_text(text, parse_mode="Markdown")
        else:
            await update.message.reply_text(f"❌ Ошибка при расчёте: {response.status_code}")

    except ValueError:
        await update.message.reply_text("❌ Пожалуйста, введи целое число.")



def format_split_result(data: dict) -> str:
    message = "📊 *Результат деления чека:*\n"

    total_amount = 0
    for i, (person, amount) in enumerate(data.items(), start=1):
        rub = amount / 100
        message += f" 👤  *{rub:.2f} ₽ с человека*\n"
        total_amount += rub

    message += "\n💬 Можешь просто скопировать этот текст и отправить друзьям:\n"
    message += "```"
    message += "\n"
    message += "\n".join([f"Отлично сегодня посидели, с тебя {amount / 100:.2f} ₽" for person, amount in data.items()])
    message += "\n📱 Скинуть можно сюда: \n"
    message += "```"

    message += "Спасибо, что пользуетесь UnBill 💸"

    return message



def format_receipt_message(receipt_data: dict) -> str:
    """Форматирует данные чека в читаемое сообщение"""
    # Проверяем наличие всех необходимых полей
    if not all(key in receipt_data for key in ['filename', 'product_units', 'total_price']):
        raise ValueError("Неполные данные в ответе сервера")
    
    message = f"📄 *Чек: {receipt_data['filename']}*\n\n"

    message += "🛍 *Товары:*\n"
    for product in receipt_data['product_units']:
        message += (
            f"  - *{product['name']}*\n"
            f"    Количество: {product['quantity']} × {product['price']/100:.2f} ₽\n"
            f"    Итого: {product['quantity'] * product['price']/100:.2f} ₽\n\n")

    message += f"💵 *Общая сумма: {receipt_data['total_price']/100:.2f} ₽*"
    
    return message

def main():
    application = Application.builder().token(BOT_TOKEN).build()
    application.add_handler(CommandHandler("start", start))
    application.add_handler(MessageHandler(filters.PHOTO, handle_photo))
    application.add_handler(CallbackQueryHandler(handle_split_choice))
    application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, handle_participant_count))


    application.run_polling()


if __name__ == "__main__":
    main()