import os
import json
import datetime

def get_data(config):
    holidays_file = config.get('holidays_file', 'holidays.json')
    today = datetime.date.today()
    results = []
    
    # Holidays
    holiday_str = "No upcoming holidays found"
    try:
        script_dir = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
        json_path = os.path.join(script_dir, holidays_file)
        if os.path.exists(json_path):
            with open(json_path, 'r') as f:
                holidays_data = json.load(f)
                
            upcoming = []
            for date_str, name in holidays_data.items():
                h_date = datetime.datetime.strptime(date_str, "%Y-%m-%d").date()
                if h_date >= today:
                    upcoming.append((h_date, name))
            
            if upcoming:
                upcoming.sort()
                next_h_date, next_h_name = upcoming[0]
                days_away = (next_h_date - today).days
                date_display = next_h_date.strftime("%d %b")
                holiday_str = f"{next_h_name} ({date_display}, {days_away} days)" if days_away > 0 else f"TODAY: {next_h_name}"
        else:
            holiday_str = f"Holidays file not found: {holidays_file}"
    except Exception as e:
        holiday_str = f"Holidays error: {e}"

    results.append(f"Nearest Holiday: {holiday_str}")

    # Sunday
    days_to_sunday = (6 - today.weekday()) % 7
    days_to_sunday = 7 if days_to_sunday == 0 else days_to_sunday
    sunday_date = (today + datetime.timedelta(days=days_to_sunday)).strftime("%d %b")
    results.append(f"Next Sunday: {sunday_date} ({days_to_sunday} days left)")
    
    return results
