import os
import pandas as pd
from PIL import Image
import json

from pathlib import Path

# 设置图片文件夹路径
IMAGE_BASE_FOLDER = 'images'


def check_folder_structure(folder_path):
    base_path = Path(folder_path)
    # 检查顶层文件夹
    for item in base_path.iterdir():
        if item.is_file():
            raise ValueError(f"当前文件夹 '{base_path}' 下存在文件 '{item.name}'，不符合要求。")
    # 检查子文件夹
    for subdir in base_path.iterdir():
        if subdir.is_dir():
            for subitem in subdir.iterdir():
                if subitem.is_dir():
                    raise ValueError(f"子文件夹 '{subdir}' 下存在文件夹 '{subitem.name}'，不符合要求。")

# 使用示例
try:
    check_folder_structure(IMAGE_BASE_FOLDER)
    print("文件夹结构符合要求。")
except ValueError as e:
    print(e)


# 初始化列表来存储文件信息
data = []

# 遍历图片文件夹及其子文件夹
for subdir, dirs, files in os.walk(IMAGE_BASE_FOLDER):
    for file in files:
        filepath = os.path.join(subdir, file)
        try:
            with Image.open(filepath) as img:
                width, height = img.size
            # 获取子文件夹路径
            subfolder = os.path.relpath(subdir, IMAGE_BASE_FOLDER)
            # 获取文件后缀
            suffix = os.path.splitext(file)[1].lower()
            # 将信息添加到数据列表中
            data.append([subfolder, file, suffix, width, height])
        except Exception as e:
            print(f"无法处理文件 {filepath}: {e}")

# 将数据列表转换为 DataFrame
df = pd.DataFrame(data, columns=['sub', 'file', 'suf', 'width', 'height'])

# 打印 DataFrame
print(df)

# 将 DataFrame 转换为字典
data_dict = df.to_dict(orient='records')

# 将字典保存为 JSON 文件
with open('images_info.json', 'w', encoding='utf-8') as f:
    json.dump(data_dict, f, ensure_ascii=False, indent=4)
