import os
import datetime
import io
import qrcode
from reportlab.lib.pagesizes import A4
from reportlab.pdfgen import canvas
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.pdfbase import pdfmetrics
from reportlab.lib.utils import ImageReader

class PDFGenerator:
    def __init__(self, config):
        self.config = config
        self.user_config = config.get('user', {})
        self.settings = config.get('settings', {})
        self.script_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
        
        self.selected_font = self.settings.get('font', 'Courier')
        self.selected_bold = self.settings.get('font_bold', 'Courier-Bold')
        self._register_fonts()

    def _register_fonts(self):
        font_paths = [
            # Linux
            "/usr/share/fonts/TTF/JetBrainsMono-Regular.ttf",
            "/usr/share/fonts/truetype/jetbrains-mono/JetBrainsMono-Regular.ttf",
            # Windows
            os.path.join(os.environ.get('WINDIR', 'C:\\Windows'), 'Fonts', 'JetBrainsMono-Regular.ttf'),
            os.path.join(os.environ.get('LOCALAPPDATA', ''), 'Microsoft', 'Windows', 'Fonts', 'JetBrainsMono-Regular.ttf'),
            # macOS
            "/Library/Fonts/JetBrainsMono-Regular.ttf",
            os.path.expanduser("~/Library/Fonts/JetBrainsMono-Regular.ttf"),
        ]
        
        bold_paths = [
            # Linux
            "/usr/share/fonts/TTF/JetBrainsMono-Bold.ttf",
            "/usr/share/fonts/truetype/jetbrains-mono/JetBrainsMono-Bold.ttf",
            # Windows
            os.path.join(os.environ.get('WINDIR', 'C:\\Windows'), 'Fonts', 'JetBrainsMono-Bold.ttf'),
            os.path.join(os.environ.get('LOCALAPPDATA', ''), 'Microsoft', 'Windows', 'Fonts', 'JetBrainsMono-Bold.ttf'),
            # macOS
            "/Library/Fonts/JetBrainsMono-Bold.ttf",
            os.path.expanduser("~/Library/Fonts/JetBrainsMono-Bold.ttf"),
        ]
        
        for fp in font_paths:
            if os.path.exists(fp):
                try:
                    pdfmetrics.registerFont(TTFont('Mono', fp))
                    self.selected_font = "Mono"
                    break
                except Exception:
                    continue
                
        for bp in bold_paths:
            if os.path.exists(bp):
                try:
                    pdfmetrics.registerFont(TTFont('MonoBold', bp))
                    self.selected_bold = "MonoBold"
                    break
                except Exception:
                    continue

    def get_qr_code(self, data):
        try:
            qr = qrcode.QRCode(
                version=1,
                error_correction=qrcode.constants.ERROR_CORRECT_L,
                box_size=10,
                border=0,
            )
            qr.add_data(data)
            qr.make(fit=True)

            img = qr.make_image(fill_color="black", back_color="white")
            img_byte_arr = io.BytesIO()
            img.save(img_byte_arr, format='PNG')
            img_byte_arr.seek(0)
            return img_byte_arr
        except Exception:
            pass
        return None

    def _clean_text(self, text):
        if not text:
            return ""
        
        import unicodedata
        
        replacements = {
            '\u2018': "'", '\u2019': "'",
            '\u201c': '"', '\u201d': '"',
            '\u2013': '-', '\u2014': '-',
            '\u2026': '...',
        }
        
        cleaned = str(text)
        for char, replacement in replacements.items():
            cleaned = cleaned.replace(char, replacement)
            
        normalized = unicodedata.normalize('NFKD', cleaned)
        result = "".join([c for c in normalized if not unicodedata.combining(c)])
        
        try:
            return result.encode('latin-1', 'replace').decode('latin-1').replace('?', ' ')
        except Exception:
            return "".join(c for c in result if ord(c) < 128)

    def generate(self, report_data, output_path, lyric=""):
        c = canvas.Canvas(output_path, pagesize=A4)
        c.setTitle(self._clean_text(f"DPPR - {datetime.date.today()}"))
        
        qr_data = self.settings.get('qr_data', "https://example.com")
        qr_stream = self.get_qr_code(qr_data)
        if qr_stream:
            qr_img = ImageReader(qr_stream)
            c.drawImage(qr_img, 485, 750, width=60, height=60, mask='auto')

        y = 800
        # Header
        name = self._clean_text(self.user_config.get('name', 'User'))
        greeting = self._clean_text(self.user_config.get('greeting', 'Good Day'))
        
        c.setFont(self.selected_bold, 12)
        c.drawString(50, y, f"{greeting}, {name}")
        y -= 20
        c.setFont(self.selected_font, 10)
        c.drawString(50, y, self._clean_text(f"DPPR (Daily Personal Printed Report): {datetime.date.today()}"))
        y -= 10
        c.drawString(50, y, "=" * 50)
        y -= 20

        # Body
        for section in report_data:
            title = section.get('title')
            lines = section.get('lines', [])
            
            if title:
                c.setFont(self.selected_bold, 10)
                c.drawString(50, y, self._clean_text(title).upper())
                y -= 15
            
            from reportlab.lib.utils import simpleSplit
            c.setFont(self.selected_font, 8)
            for line in lines:
                wrapped_lines = simpleSplit(str(line), self.selected_font, 8, 500)
                for wrapped_line in wrapped_lines:
                    cleaned_line = self._clean_text(wrapped_line)
                    c.drawString(50, y, cleaned_line)
                    y -= 10
                    
                    if y < 80: 
                        c.showPage()
                        y = 800
                        c.setFont(self.selected_font, 8)
            
            y -= 10 
            if y < 80:
                c.showPage()
                y = 800

        if lyric:
            from reportlab.lib.utils import simpleSplit
            c.setFont(self.selected_font, 7)
            
            max_width = 380 
            lines = simpleSplit(self._clean_text(lyric), self.selected_font, 7, max_width)
            
            lyric_y = 70
            for line in lines:
                c.drawString(50, lyric_y, line)
                lyric_y -= 8

        # Signature
        sig_path = self.settings.get('signature_path', "signature.png")
        sig_abs_path = os.path.join(self.script_dir, sig_path)
        if os.path.exists(sig_abs_path):
            c.setFont(self.selected_bold, 8)
            c.drawRightString(545, 60, "Verified and Authorised:")
            c.drawImage(sig_abs_path, 445, 20, width=100, height=40, mask='auto', preserveAspectRatio=True)
            
        c.save()
        return output_path
