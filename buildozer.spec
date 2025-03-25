[app]
title = Unbill
package.name = unbill
package.domain = org.test

# Source directory and included files
source.dir = src
source.include_exts = py,png,jpg,kv,atlas

# Application version
version = 0.1

# Required dependencies
requirements = python3,kivy,plyer

# Application orientation
orientation = portrait

# Fullscreen setting
fullscreen = 0

# Android-specific configurations
android.api = 31
android.minapi = 21
android.ndk = 25b
android.ndk_api = 24
android.archs = arm64-v8a, armeabi-v7a

# Permissions required by Plyer filechooser
android.permissions = READ_EXTERNAL_STORAGE, WRITE_EXTERNAL_STORAGE

# Enable Android backup
android.allow_backup = True

# Output format
android.debug_artifact = apk
android.release_artifact = aab

[buildozer]
log_level = 2
warn_on_root = 1
