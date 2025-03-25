from kivy.app import App
from kivy.uix.boxlayout import BoxLayout
from kivy.uix.button import Button
from kivy.uix.image import Image
from kivy.uix.filechooser import FileChooserIconView
from plyer import filechooser

class ImagePickerApp(App):
    def build(self):
        self.layout = BoxLayout(orientation='vertical')
        
        # Button to open file chooser
        self.choose_button = Button(text='Выбрать изображение', size_hint=(1, 0.1))
        self.choose_button.bind(on_press=self.pick_image)
        
        # Image widget to display selected image
        self.img = Image(size_hint=(1, 0.9))
        
        self.layout.add_widget(self.choose_button)
        self.layout.add_widget(self.img)
        
        return self.layout
    
    def pick_image(self, instance):
        filechooser.open_file(on_selection=self.selected)
    
    def selected(self, selection):
        if selection:
            self.img.source = selection[0]  # Set image source to the selected file
            self.img.reload()

if __name__ == '__main__':
    ImagePickerApp().run()