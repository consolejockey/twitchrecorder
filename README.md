# twitch-recorder

This is a lightweight recorder written in Go that utilizes [Streamlink](https://github.com/streamlink/streamlink) to record live streams from Twitch. Streamlink is a command-line utility that extracts streams from various services. twitch-recorder integrates Streamlink with a live stream status check, enabling automated Twitch stream recording. With this functionality, users can effortlessly record their desired streams without manual intervention, streamlining the entire recording process.

## Usage
1. Ensure Streamlink is installed on your system.
2. Create a Twitch application to obtain a client ID and client secret. You can do this by visiting the [Twitch Developer Dashboard](https://dev.twitch.tv/) and registering a new application.
3. Configure the config.json file with the obtained client ID, client secret, Twitch streamer's name, preferred quality, and output location.
4. Run the recorder.
5. When your configured streamer goes live, a .mp4 file with a timestamp will be created.
