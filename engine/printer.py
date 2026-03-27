import os

def print_pdf(pdf_path, config):
    try:
        import cups
    except ImportError:
        print("Error: `pycups` not installed. Cannot print.")
        return False

    settings = config.get('settings', {})
    printing_config = settings.get('printing', {})
    
    if not printing_config.get('enabled', False):
        return False

    try:
        conn = cups.Connection()
        printers = conn.getPrinters()
        if not printers:
            print("No printers configured in CUPS.")
            return False
        
        printer_name = printing_config.get('cups_printer') or settings.get('cups_printer')
        if not printer_name:
            printer_name = list(printers.keys())[0]
        
        if printer_name not in printers:
            print(f"Printer '{printer_name}' not found. Available: {', '.join(printers.keys())}")
            return False

        conn.printFile(printer_name, pdf_path, "Daily Automation Report", {})
        print(f"Sent to printer: {printer_name}")
        return True
    except Exception as e:
        print(f"Printing failed: {e}")
        return False
