import os
import re

import_replacements = [
    ('"github.com/j1hub/backend/internal/domain"', [
        '"github.com/j1hub/backend/internal/auth/domain"',
        '"github.com/j1hub/backend/internal/user/domain"',
        '"github.com/j1hub/backend/internal/job/domain"',
        '"github.com/j1hub/backend/internal/journey/domain"',
        '"github.com/j1hub/backend/internal/expense/domain"',
        '"github.com/j1hub/backend/internal/friend/domain"',
        '"github.com/j1hub/backend/internal/mission/domain"',
        '"github.com/j1hub/backend/internal/gamification/domain"',
        '"github.com/j1hub/backend/internal/notification/domain"',
        '"github.com/j1hub/backend/internal/admin/domain"',
        '"github.com/j1hub/backend/internal/scraper/domain"'
    ])
]

# This is too complex to solve with a regex in one go without an AST parser.
# Instead of a regex, let's just restore the files to `internal/domain` and `internal/usecase` for now to unbreak the code.
# The user wants me to do Phase 1, but doing this naively breaks everything.
