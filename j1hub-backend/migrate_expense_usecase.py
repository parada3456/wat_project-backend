import os
import re

ROOT_DIR = "/Users/user/development/work/WAT_project/backend/j1hub-backend"
INTERNAL_DIR = os.path.join(ROOT_DIR, "internal")
EXPENSE_DIR = os.path.join(INTERNAL_DIR, "expense")
USECASE_DIR = os.path.join(INTERNAL_DIR, "usecase")

def migrate_expense_usecase():
    # 1. Move files
    old_uc = os.path.join(USECASE_DIR, "manage_expense.go")
    new_uc = os.path.join(EXPENSE_DIR, "usecase", "manage_expense.go")
    old_test = os.path.join(USECASE_DIR, "manage_expense_test.go")
    new_test = os.path.join(EXPENSE_DIR, "usecase", "manage_expense_test.go")
    
    for old_path, new_path in [(old_uc, new_uc), (old_test, new_test)]:
        if os.path.exists(old_path):
            with open(old_path, "r") as f:
                content = f.read()
            content = content.replace("package usecase", "package expenseusecase")
            with open(new_path, "w") as f:
                f.write(content)
            os.remove(old_path)
            print(f"Moved {os.path.basename(old_path)}")

    # 2. Update handler and main.go
    handler_path = os.path.join(INTERNAL_DIR, "adapter", "http", "handler", "expense_handler.go")
    main_path = os.path.join(ROOT_DIR, "cmd", "server", "main.go")
    
    for filepath in [handler_path, main_path]:
        if not os.path.exists(filepath): continue
        with open(filepath, "r") as f:
            content = f.read()
            
        orig_content = content
        
        # Replace the function calls and types
        content = content.replace("usecase.ManageExpenseUseCase", "expenseusecase.ManageExpenseUseCase")
        content = content.replace("usecase.NewManageExpenseUseCase", "expenseusecase.NewManageExpenseUseCase")
        
        if "expenseusecase." in content and '"github.com/j1hub/backend/internal/expense/usecase"' not in content:
            if '"github.com/j1hub/backend/internal/usecase"' in content:
                content = content.replace(
                    '"github.com/j1hub/backend/internal/usecase"',
                    '"github.com/j1hub/backend/internal/usecase"\n\texpenseusecase "github.com/j1hub/backend/internal/expense/usecase"'
                )
            else:
                content = re.sub(
                    r'import \(\n',
                    'import (\n\texpenseusecase "github.com/j1hub/backend/internal/expense/usecase"\n',
                    content, count=1
                )
        
        if content != orig_content:
            with open(filepath, "w") as f:
                f.write(content)
            print(f"Updated {os.path.basename(filepath)}")

if __name__ == "__main__":
    migrate_expense_usecase()
