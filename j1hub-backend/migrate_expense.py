import os
import re

ROOT_DIR = "/Users/user/development/work/WAT_project/backend/j1hub-backend"
INTERNAL_DIR = os.path.join(ROOT_DIR, "internal")

expense_types = [
    "ExpenseTransaction", "ExpenseSplit", "PaymentStatus", "ApprovalStatus",
    "PaymentPending", "PaymentSubmitted", "PaymentApproved", "PaymentOverdue",
    "ApprovalPending", "ApprovalApproved", "ApprovalRejected"
]

def migrate_expense_feature():
    # 1. Create directories
    os.makedirs(os.path.join(INTERNAL_DIR, "expense", "domain"), exist_ok=True)
    os.makedirs(os.path.join(INTERNAL_DIR, "expense", "usecase"), exist_ok=True)
    os.makedirs(os.path.join(INTERNAL_DIR, "expense", "port"), exist_ok=True)
    os.makedirs(os.path.join(INTERNAL_DIR, "expense", "adapter", "postgres"), exist_ok=True)
    os.makedirs(os.path.join(INTERNAL_DIR, "expense", "adapter", "http"), exist_ok=True)

    # 2. Move domain/expense.go to expense/domain/expense.go
    old_domain = os.path.join(INTERNAL_DIR, "domain", "expense.go")
    new_domain = os.path.join(INTERNAL_DIR, "expense", "domain", "expense.go")
    
    if os.path.exists(old_domain):
        with open(old_domain, "r") as f:
            content = f.read()
        content = content.replace("package domain", "package expensedomain")
        with open(new_domain, "w") as f:
            f.write(content)
        os.remove(old_domain)
        print("Moved expense.go")
        
    # 3. Update all references
    for root, dirs, files in os.walk(ROOT_DIR):
        for file in files:
            if not file.endswith(".go"): continue
            
            filepath = os.path.join(root, file)
            with open(filepath, "r") as f:
                content = f.read()
                
            orig_content = content
            
            for t in expense_types:
                content = re.sub(r'\bdomain\.' + t + r'\b', f'expensedomain.{t}', content)
            
            if "expensedomain." in content and '"github.com/j1hub/backend/internal/expense/domain"' not in content:
                if '"github.com/j1hub/backend/internal/domain"' in content:
                    content = content.replace(
                        '"github.com/j1hub/backend/internal/domain"',
                        '"github.com/j1hub/backend/internal/domain"\n\texpensedomain "github.com/j1hub/backend/internal/expense/domain"'
                    )
                else:
                    content = re.sub(
                        r'import \(\n',
                        'import (\n\texpensedomain "github.com/j1hub/backend/internal/expense/domain"\n',
                        content, count=1
                    )
            
            if content != orig_content:
                with open(filepath, "w") as f:
                    f.write(content)

if __name__ == "__main__":
    migrate_expense_feature()
