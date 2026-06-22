import os
import re

ROOT_DIR = "/Users/user/development/work/WAT_project/backend/j1hub-backend"
INTERNAL_DIR = os.path.join(ROOT_DIR, "internal")

# Feature mapping
# Feature -> (domain file, usecase files, domain types to replace)
features = {
    "user": ("user.go", ["manage_user.go", "manage_user_test.go"], ["User", "Profile", "UserRole", "UserStatus", "ProfileVisibility", "CreditScore"]),
    "auth": (None, ["login_user.go", "login_user_test.go", "register_user.go", "register_user_test.go"], []),
    "friend": ("friendship.go", ["manage_friendship.go", "manage_friendship_test.go"], ["Friendship", "FriendshipStatus"]),
    "mission": ("mission.go", ["manage_mission.go", "manage_mission_test.go", "complete_mission.go", "complete_mission_test.go"], ["Mission", "UserMission", "UserMissionStatus", "Task", "UserTask", "UserTaskStatus", "PointReward"]),
    "gamification": ("gamification.go", ["leaderboard.go", "leaderboard_test.go", "radar.go", "radar_test.go", "reward_engine.go", "reward_engine_test.go"], ["JourneyPhase", "UserPhaseHistory", "PointLedger", "Badge", "UserBadge", "TriggerType"]),
    "notification": ("notification.go", ["manage_notification.go", "manage_notification_test.go"], ["Notification", "NotificationType"]),
    "admin": (None, ["admin_usecase.go"], [])
}

def migrate_all():
    for feature, (domain_file, usecase_files, types) in features.items():
        # 1. Create directories
        for layer in ["domain", "usecase", "port", "adapter/http", "adapter/postgres"]:
            os.makedirs(os.path.join(INTERNAL_DIR, feature, layer), exist_ok=True)
            
        # 2. Move domain file
        if domain_file:
            old_domain = os.path.join(INTERNAL_DIR, "domain", domain_file)
            new_domain = os.path.join(INTERNAL_DIR, feature, "domain", domain_file)
            if os.path.exists(old_domain):
                with open(old_domain, "r") as f:
                    content = f.read()
                content = content.replace("package domain", f"package {feature}domain")
                with open(new_domain, "w") as f:
                    f.write(content)
                os.remove(old_domain)
                print(f"Moved {domain_file} to {feature}/domain")
                
        # 3. Move usecase files
        for uc_file in usecase_files:
            old_uc = os.path.join(INTERNAL_DIR, "usecase", uc_file)
            new_uc = os.path.join(INTERNAL_DIR, feature, "usecase", uc_file)
            if os.path.exists(old_uc):
                with open(old_uc, "r") as f:
                    content = f.read()
                content = content.replace("package usecase", f"package {feature}usecase")
                with open(new_uc, "w") as f:
                    f.write(content)
                os.remove(old_uc)
                print(f"Moved {uc_file} to {feature}/usecase")

        # 4. Global string replacements
        for root_path, dirs, files in os.walk(ROOT_DIR):
            for file in files:
                if not file.endswith(".go"): continue
                filepath = os.path.join(root_path, file)
                
                with open(filepath, "r") as f:
                    content = f.read()
                    
                orig_content = content
                
                # Replace domain types
                for t in types:
                    content = re.sub(r'\bdomain\.' + t + r'\b', f'{feature}domain.{t}', content)
                    
                # Replace usecase interface/struct instantiation
                for uc_file in usecase_files:
                    uc_base = uc_file.replace('.go', '').replace('_test', '')
                    # naive camelcase heuristic: manage_user -> ManageUser
                    camel = ''.join(word.title() for word in uc_base.split('_'))
                    content = content.replace(f"usecase.{camel}UseCase", f"{feature}usecase.{camel}UseCase")
                    content = content.replace(f"usecase.New{camel}UseCase", f"{feature}usecase.New{camel}UseCase")
                    content = content.replace(f"usecase.{camel}Cmd", f"{feature}usecase.{camel}Cmd")
                
                # Add import if replaced
                if f"{feature}domain." in content and f'"github.com/j1hub/backend/internal/{feature}/domain"' not in content:
                    content = re.sub(r'import \(\n', f'import (\n\t{feature}domain "github.com/j1hub/backend/internal/{feature}/domain"\n', content, count=1)
                if f"{feature}usecase." in content and f'"github.com/j1hub/backend/internal/{feature}/usecase"' not in content:
                    content = re.sub(r'import \(\n', f'import (\n\t{feature}usecase "github.com/j1hub/backend/internal/{feature}/usecase"\n', content, count=1)

                if content != orig_content:
                    with open(filepath, "w") as f:
                        f.write(content)

if __name__ == "__main__":
    migrate_all()
