import os
import re
import shutil

ROOT_DIR = "/Users/user/development/work/WAT_project/backend/j1hub-backend"
INTERNAL_DIR = os.path.join(ROOT_DIR, "internal")

# Targets to replace
job_types = [
    "JobPosting", "JobHousing", "JobOverallRating", "JobReview", 
    "UserCart", "CartStatus", "CartSaved", "CartViewed", "CartApplied", "CartRemoved"
]

def migrate_job_feature():
    # 1. Move domain/job.go to job/domain/job.go
    old_domain = os.path.join(INTERNAL_DIR, "domain", "job.go")
    new_domain = os.path.join(INTERNAL_DIR, "job", "domain", "job.go")
    
    if os.path.exists(old_domain):
        with open(old_domain, "r") as f:
            content = f.read()
        
        # Change package name
        content = content.replace("package domain", "package jobdomain")
        
        with open(new_domain, "w") as f:
            f.write(content)
        
        os.remove(old_domain)
        print("Moved and updated job.go")
        
    # 2. Update all references in the codebase
    for root, dirs, files in os.walk(ROOT_DIR):
        for file in files:
            if not file.endswith(".go"):
                continue
            
            filepath = os.path.join(root, file)
            with open(filepath, "r") as f:
                content = f.read()
                
            orig_content = content
            
            # Replace domain.<Type> with jobdomain.<Type>
            for t in job_types:
                content = re.sub(r'\bdomain\.' + t + r'\b', f'jobdomain.{t}', content)
            
            # If jobdomain is used, ensure the import is present
            if "jobdomain." in content and '"github.com/j1hub/backend/internal/job/domain"' not in content:
                # Add import using a simple heuristic: find existing domain import and add next to it
                if '"github.com/j1hub/backend/internal/domain"' in content:
                    content = content.replace(
                        '"github.com/j1hub/backend/internal/domain"',
                        '"github.com/j1hub/backend/internal/domain"\n\tjobdomain "github.com/j1hub/backend/internal/job/domain"'
                    )
                else:
                    # Just an approximation for files that might not have domain imported
                    content = re.sub(
                        r'import \(\n',
                        'import (\n\tjobdomain "github.com/j1hub/backend/internal/job/domain"\n',
                        content, count=1
                    )
            
            if content != orig_content:
                with open(filepath, "w") as f:
                    f.write(content)
                print(f"Updated references in {filepath}")

if __name__ == "__main__":
    migrate_job_feature()
