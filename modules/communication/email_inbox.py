import os
import imaplib
import email
from email.header import decode_header

def get_data(config):
    
    host = os.getenv('EMAIL_HOST', 'imap.gmail.com')
    user = os.getenv('EMAIL_USER')
    password = os.getenv('EMAIL_PASS')
    
    if not user or not password:
        return ["Email Error: Missing EMAIL_USER or EMAIL_PASS in .env"]
        
    count = config.get('count', 3)
    results = []
    
    try:
        mail = imaplib.IMAP4_SSL(host)
        mail.login(user, password)
        mail.select("inbox")
        
        status, messages = mail.search(None, 'UNSEEN')
        msg_ids = messages[0].split()
        
        if not msg_ids:
            status, messages = mail.search(None, 'ALL')
            msg_ids = messages[0].split()
        
        for msg_id in reversed(msg_ids[-count:]):
            res, msg_data = mail.fetch(msg_id, "(RFC822)")
            for response_part in msg_data:
                if isinstance(response_part, tuple):
                    msg = email.message_from_bytes(response_part[1])
                    
                    subject, encoding = decode_header(msg["Subject"])[0]
                    if isinstance(subject, bytes):
                        subject = subject.decode(encoding if encoding else "utf-8")
                        
                    from_, encoding = decode_header(msg.get("From"))[0]
                    if isinstance(from_, bytes):
                        from_ = from_.decode(encoding if encoding else "utf-8")
                    
                    sender = from_.split('<')[0].strip()
                    
                    results.append(f"{sender}: {subject}")
                    
        mail.close()
        mail.logout()
        
        if not results:
            results.append("Inbox is empty.")
            
    except Exception as e:
        results.append(f"IMAP Error: {e}")
        
    return results
