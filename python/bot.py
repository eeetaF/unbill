from telegram import InlineKeyboardButton, InlineKeyboardMarkup, Update
from telegram.ext import ApplicationBuilder, CommandHandler, CallbackQueryHandler, ContextTypes
from config import BOT_TOKEN

async def start_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    



# Команда /menu показывает клавиатуру
async def menu_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    keyboard = [
        [
            InlineKeyboardButton("Разделить чек", callback_data="split"),
            InlineKeyboardButton("Помощь", callback_data="help")
        ]
    ]
    reply_markup = InlineKeyboardMarkup(keyboard)
    await update.message.reply_text("Выбери действие:", reply_markup=reply_markup)



async def button_handler(update: Update, context: ContextTypes.DEFAULT_TYPE):
    query = update.callback_query
    await query.answer() 

    if query.data == "split":
        await query.edit_message_text("🧾 Начнём делить чек! Отправь фото или сумму.")
    elif query.data == "help":
        await query.edit_message_text("ℹ️ Я помогу тебе разделить счёт по чеку. Просто отправь его!")



def main():
    app = ApplicationBuilder().token(BOT_TOKEN).build()

    app.add_handler(CommandHandler("menu", menu_command))
    app.add_handler(CallbackQueryHandler(button_handler))

    print("Бот запущен с клавиатурой!")
    app.run_polling()


if __name__ == "__main__":
    main()
