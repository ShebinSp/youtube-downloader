import sys
from moviepy.editor import VideoFileClip, AudioFileClip

def merge(video_path, audio_path, output_path):
    try:
        video = VideoFileClip(video_path)
        audio = AudioFileClip(audio_path)
        final_video = video.set_audio(audio)
        final_video.write_videofile(output_path, codec="libx264", audio_codec="aac")
    except Exception as e:
        print(f"Error during merging: {e}")
        

if __name__ == "__main__":
    merge(sys.argv[1], sys.argv[2], sys.argv[3])
