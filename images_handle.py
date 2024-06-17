import hashlib
import shutil
from pathlib import Path

def compute_md5(file_path):
    """计算文件的 MD5 值"""
    hasher = hashlib.md5()
    with open(file_path, 'rb') as f:
        buffer = f.read(65536)  # 64KB 缓冲区
        while len(buffer) > 0:
            hasher.update(buffer)
            buffer = f.read(65536)
    return hasher.hexdigest()

def main(source_dir, target_dir):
    """主函数：计算文件的 MD5 值，重命名并处理重复文件"""
    source_path = Path(source_dir)
    target_path = Path(target_dir)
    
    # 确保目标文件夹存在
    target_path.mkdir(parents=True, exist_ok=True)

    # 用于存储已处理的文件的 MD5 值
    md5_set = set()

    # 遍历源文件夹中的文件
    for file_path in source_path.iterdir():
        if file_path.is_file():
            # 计算文件的 MD5 值
            md5 = compute_md5(file_path)
            
            # 构造新文件名，保留文件后缀
            new_filename = f"{md5}{file_path.suffix}"
            target_file_path = target_path / new_filename

            if md5 in md5_set:
                # 如果 MD5 值已存在于集合中，说明是重复文件，删除源文件
                print(f"删除重复文件: {file_path}")
                file_path.unlink()
            else:
                # 将文件移动到目标文件夹并重命名
                shutil.move(str(file_path), str(target_file_path))
                print(f"移动文件: {file_path} 到 {target_file_path}")

                # 将当前文件的 MD5 值添加到集合中
                md5_set.add(md5)

if __name__ == "__main__":
    source_dir = "images/dongman"
    target_dir = "images/dongman"
    main(source_dir, target_dir)