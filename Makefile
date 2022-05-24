launch: dial_habit_tracker.apk
	if ! sdk/platform-tools/adb wait-for-local-device; then \
		echo "maybe try make start?" >&2; \
		exit 1; \
	fi
	sdk/platform-tools/adb install -r dial_habit_tracker.apk 
	sdk/platform-tools/adb shell am start -n us.cernu.DialHabitTracker/org.golang.app.GoNativeActivity

start: sdk/emulator/emulator ~/.android/avd/dial-habit-tracker.avd
	sdk/emulator/emulator \
		-accel on \
		@dial-habit-tracker &

dial_habit_tracker.apk: dep go.mod main.go sdk/ndk/21.4.7075529/ndk-build AndroidManifest.xml Icon.png
	go install fyne.io/fyne/v2/cmd/fyne@latest
	ANDROID_NDK_HOME="$$(readlink -f sdk/ndk/21.4.7075529)" \
	fyne package -os android -appID us.cernu.DialHabitTracker

dep: .dep
.dep:
	sudo apt-get install libc-dev libegl1-mesa-dev libgles2-mesa-dev libx11-dev
	touch .dep

sdk/cmdline-tools:
	mkdir -p sdk/cmdline-tools
	cd sdk/cmdline-tools; \
	wget https://dl.google.com/android/repository/commandlinetools-linux-8512546_latest.zip; \
	unzip commandlinetools-linux-8512546_latest.zip; \
	rm commandlinetools-linux-8512546_latest.zip; \
	mv cmdline-tools tools

sdk/platforms/android-21: sdk/cmdline-tools
	sdk/cmdline-tools/tools/bin/sdkmanager --sdk_root=sdk platforms\;android-21

sdk/ndk/21.4.7075529/ndk-build: sdk/cmdline-tools
	sdk/cmdline-tools/tools/bin/sdkmanager --sdk_root=sdk ndk\;21.4.7075529
	touch -c sdk/ndk/21.4.7075529/ndk-build

sdk/emulator/emulator: sdk/cmdline-tools
	sdk/cmdline-tools/tools/bin/sdkmanager --sdk_root=sdk emulator

sdk/system-images/android-21/default/x86_64: sdk/cmdline-tools
	sdk/cmdline-tools/tools/bin/sdkmanager --sdk_root=sdk 'system-images;android-21;default;x86_64'

~/.android/avd/dial-habit-tracker.avd:
	cd sdk; \
	export ANDROID_HOME="$$(readlink -f .)"; \
	export ANDROID_SDK_ROOT="$$(readlink -f .)"; \
	export ANDROID_AVD_HOME="~/.android/avd"; \
	echo no | \
	./cmdline-tools/tools/bin/avdmanager \
		--verbose \
		--clear-cache \
		create avd \
		--force \
		-k 'system-images;android-21;default;x86_64' \
		--name dial-habit-tracker
