import json
import os
import requests

# Input JSON file with champion data
#  TODO: fix this to read from champion portion of the JSON
JSON_PATH = 'champions.json'

# Output directory for images
OUTPUT_DIR = './frontend/public/tft-champion-icons'
os.makedirs(OUTPUT_DIR, exist_ok=True)

# CommunityDragon PNG base URL
BASE_URL = 'https://raw.communitydragon.org/latest/game/assets/'
# doc for asset mapping https://communitydragon.org/documentation/assets

def tex_to_png_url(tex_path):
    """
    Convert a .tex asset path to a CommunityDragon .png URL.
    """
    if not tex_path.lower().endswith('.tex'):
        return None
    
    fileName = tex_path.split('/')[-1]
    prefix = fileName.split('.')[0] + '_square'
    suffix = fileName.split('.')[1:-1]
    print("suffix: ", suffix)
    download_file = '.'.join([prefix] + suffix).lower() + ".png"
    print("fileName: ", download_file)
    directory = '/'.join(tex_path.split('/')[1:-1])
    print("directory: ", directory)
    
    relative_path = directory.lower().replace('skins/base/images', 'hud/')
    print("relative_path: ", relative_path)
    png_path = relative_path + download_file
    return BASE_URL + png_path

def download_image(url, output_path):
    """
    Download an image from URL and save it to output_path.
    """
    try:
        response = requests.get(url)
        response.raise_for_status()
        with open(output_path, 'wb') as f:
            f.write(response.content)
        print(f"✅ Downloaded: {output_path} from {url}")
    except requests.RequestException as e:
        print(f"❌ Failed to download {url}: {e}")

def main():
    with open(JSON_PATH, 'r', encoding='utf-8') as f:
        data = json.load(f)
        champions = data["champions"]

    for champ in champions:
        tex_path = champ.get('icon')
        if not tex_path:
            continue

        png_url = tex_to_png_url(tex_path)
        if not png_url:
            continue

        file_name = os.path.basename(tex_path.replace('.tex', '.png'))
        output_path = os.path.join(OUTPUT_DIR, file_name)
        print(png_url)

        download_image(png_url, output_path)

if __name__ == "__main__":
    main()
