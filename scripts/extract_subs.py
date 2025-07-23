import sys
import json
import time
from videocr import save_subtitles_to_file
import yt_dlp
# import paddle

def main():
    try:
        # Read all stdin
        input_data = sys.stdin.read()
        task = json.loads(input_data)

        # Access fields
        task_id = task.get("task_id")
        video_url = task.get("video_url")
        language_code = task.get("language_code")
        start_time = task.get("start_time")
        end_time = task.get("end_time")
        frames_to_skip = task.get("frames_to_skip")
        subtitle_box = task.get("subtitle_box", {})

        top_left = subtitle_box.get("top_left", {})
        bottom_right = subtitle_box.get("bottom_right", {})

        print(f"Task ID: {task_id}")
        print(f"Video URL: {video_url}")
        print(f"Language: {language_code}")
        print(f"Start: {start_time}, End: {end_time}")
        print(f"Frames to skip: {frames_to_skip}")
        print(f"Box: {top_left} to {bottom_right}")

        download_video_only(video_url, ".", "video-"+task_id)

        # is_gpu_available()

        # save_subtitles_to_file(
        #     input_file_path,
        #     output_file_path,
        #     lang=language_code,
        #     time_start=start_time,
        #     time_end=end_time,
        #     conf_threshold=confidence_threshold,
        #     sim_threshold=similarity_threshold,
        #     use_gpu=use_gpu,
        #     # # Models different from the default mobile models can be downloaded here: https://github.com/PaddlePaddle/PaddleOCR/blob/release/2.3/doc/doc_en/models_list_en.md
        #     # det_model_dir='<PADDLEOCR DETECTION MODEL DIR>', rec_model_dir='<PADDLEOCR RECOGNITION MODEL DIR>',
        #     # brightness_threshold=210, similar_image_threshold=1000 # filters might help
        #     # use_fullframe=True, # note: videocr just assumes horizontal lines of text. vertical text scenario hasn't been implemented yet
        #     frames_to_skip=frames_to_skip, # can skip inference for some frames to speed up the process
        #     crop_x=crop_x, crop_y=crop_y, crop_width=crop_width, crop_height=crop_height)

        # time.sleep(5)

        # Your processing logic here...

    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


def is_gpu_available():
    """
    Checks if PaddlePaddle can use a GPU.
    Returns True if GPU is compiled and available, False otherwise.
    """
    if paddle.fluid.is_compiled_with_cuda():
        if paddle.cuda.is_available():
            print("GPU is available for PaddlePaddle.")
            print(f"Number of GPUs: {paddle.cuda.device_count()}")
            # Optionally print GPU names
            for i in range(paddle.cuda.device_count()):
                print(f"  GPU {i}: {paddle.device.get_device_name(f'gpu:{i}')}")
            return True
        else:
            print("PaddlePaddle is compiled with CUDA, but CUDA devices are not currently available.")
            return False
    else:
        print("PaddlePaddle is NOT compiled with CUDA (GPU support).")
        return False


def download_video_only(url, output_path='.', file_name=None):
    """
    Downloads only the video stream (no audio) from a given URL using yt-dlp.

    Args:
        url (str): The URL of the video (e.g., YouTube link).
        output_path (str): The directory to save the video. Defaults to current directory.
        file_name (str, optional): The desired name for the downloaded file (without extension).
                                   If None, yt-dlp will use its default naming.
    """
    ydl_opts = {
        'format': 'bestvideo', # Only download the best available video stream, without audio
        'outtmpl': f'{output_path}/{file_name}.%(ext)s' if file_name else f'{output_path}/%(title)s.%(ext)s',
        'noplaylist': True, # Ensure only single video is downloaded even if part of a playlist
        # 'merge_output_format': 'mp4', # Not strictly needed as we're not merging, but can be kept for consistency
        'quiet': False,
        'no_warnings': False,
        'progress_hooks': [lambda d: print(d['status']) if d['status'] == 'finished' else None],
    }

    try:
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            # ydl.download([url])
            info_dict = ydl.extract_info(url, download=True)
            # info_dict = ydl.extract_info(url)
            print(f"\nSuccessfully downloaded video-only: {info_dict.get('title', 'Video')}")
            # print(info_dict)
    except Exception as e:
        print(f"\nAn error occurred: {e}")






if __name__ == "__main__":
    main()
