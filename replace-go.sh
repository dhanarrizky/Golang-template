#!/bin/bash

# ===========================
# Default old text (shown to the user)
# ===========================
OLD_TEXT="github.com/dhanarrizky/Golang-template"

echo "Current old text to replace:"
echo "→ $OLD_TEXT"
echo ""

# Ask user for new replacement text
read -p "Enter the new replacement text: " NEW_TEXT

echo ""
echo "Replacing:"
echo "  '$OLD_TEXT'"
echo "     →"
echo "  '$NEW_TEXT'"
echo ""

# ===========================
# List of files to update
# ===========================
FILES=(
    "./cmd/app/main.go"
    "./pkg/auth/jwt.go"
    "./pkg/responses/response.go"
    "./pkg/utils/error.go"
    "./pkg/utils/id_generator.go"
    "./internal/config/config.go"
    "./internal/delivery/routes.go"
    "./internal/delivery/http/middleware/auth.go"
    "./internal/delivery/http/middleware/cors.go"
    "./internal/delivery/http/middleware/logging.go"
    "./internal/delivery/http/middleware/ratelimit.go"
    "./internal/delivery/http/middleware/recovery.go"
    "./internal/delivery/http/handlers/auth_handlers.go"
    "./internal/delivery/http/handlers/user_handlers.go"
    "./internal/domain/entities/user.go"
    "./internal/domain/valueobjects/email.go"
    "./internal/infrastructure/cache/redis.go"
    "./internal/infrastructure/database/postgres/db.go"
    "./internal/infrastructure/database/postgres/user_repository.go"
    "./internal/repository/user_repository.go"
    "./internal/usecase/user/create_user.go"
    "./internal/usecase/user/delete_user.go"
    "./internal/usecase/user/get_user.go"
    "./go.mod"
)

# ===========================
# Processing loop
# ===========================
for FILE_PATH in "${FILES[@]}"; do

    echo "Processing file: $FILE_PATH"

    # Check if file exists
    if [ ! -f "$FILE_PATH" ]; then
        echo "❌ File not found, skipping: $FILE_PATH"
        continue
    fi

    # Replace content inside file
    sed -i "s|${OLD_TEXT}|${NEW_TEXT}|g" "$FILE_PATH"

    echo "✅ Replacement completed: $FILE_PATH"
    echo ""
done

echo "🎉 All replacements completed!"
