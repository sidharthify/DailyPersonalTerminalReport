import os
import yaml
import importlib
import argparse
import sys
from dotenv import load_dotenv
from engine.pdf_generator import PDFGenerator
from engine.printer import print_pdf

def find_module_path(module_name, search_dir):
    for root, dirs, files in os.walk(search_dir):
        if f"{module_name}.py" in files:
            rel_path = os.path.relpath(root, search_dir)
            if rel_path == '.':
                return f"modules.{module_name}"
            dot_path = rel_path.replace(os.sep, '.')
            return f"modules.{dot_path}.{module_name}"
    return None

def main():
    parser = argparse.ArgumentParser(description="DPPR: Daily Personal Printed Report")
    parser.add_argument("--config", default="config.yaml", help="Path to config.yaml")
    parser.add_argument("--no-print", action="store_true", help="Generate PDF but do not print")
    args = parser.parse_args()

    script_dir = os.path.dirname(os.path.abspath(__file__))
    load_dotenv(os.path.join(script_dir, '.env'))

    if not os.path.exists(args.config):
        print(f"Error: Config file {args.config} not found.")
        return

    with open(args.config, 'r') as f:
        config = yaml.safe_load(f) or {}

    report_data = []
    lyric = ""
    
    settings = config.get('settings', {})
    master_curr = settings.get('master_currency')
    
    modules_dir = os.path.join(script_dir, "modules")
    layout = config.get('layout', [])
    for entry in layout:
        module_name = entry.get('module')
        section_title = entry.get('title')
        module_config = entry.get('config', {})
        
        user_config = config.get('user', {})
        settings = config.get('settings', {})
        
        master_curr = settings.get('master_currency')
        if master_curr and 'master_currency' not in module_config:
            module_config['master_currency'] = master_curr
            
        if 'lat' not in module_config and 'lat' in user_config:
            module_config['lat'] = user_config['lat']
        if 'lon' not in module_config and 'lon' in user_config:
            module_config['lon'] = user_config['lon']

        if module_name == 'news':
            module_config['feeds'] = entry.get('feeds', [])

        try:
            dot_path = find_module_path(module_name, modules_dir)
            if not dot_path:
                raise ImportError(f"Module '{module_name}' not found in {modules_dir}")

            module = importlib.import_module(dot_path)
            
            if module_name == 'quotes':
                if hasattr(module, 'get_data'):
                    lines = module.get_data(module_config)
                    if lines:
                        lyric = lines[0]
                continue
                
            if hasattr(module, 'get_data'):
                lines = module.get_data(module_config)
                report_data.append({
                    "title": section_title,
                    "lines": lines
                })
                
        except Exception as e:
            print(f"Error loading module {module_name}: {e}")
            report_data.append({
                "title": section_title or module_name,
                "lines": [f"Error loading module: {e}"]
            })

    pdf_gen = PDFGenerator(config)
    pdf_path = os.path.join(script_dir, "daily_report.pdf")
    pdf_gen.generate(report_data, pdf_path, lyric=lyric)
    print(f"✅ DPPR generated at: {pdf_path}")

    if args.no_print:
        print("Skipping print due to --no-print flag.")
        return

    settings = config.get('settings', {})
    printing_config = settings.get('printing', {})
    
    if printing_config.get('enabled', False):
        print_pdf(pdf_path, config)
    else:
        print("Skipping print (disabled in config).")

if __name__ == "__main__":
    main()
