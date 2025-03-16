from typing import List
import pandas as pd
import os
from langchain.schema import Document


class GoodsLoader:
    def __init__(self, data_path: str = "./data/goods.xlsx"):
        self.data_path = data_path

    def load_Goods(self) -> List[Document]:
        """从Excel文件加载商品数据并转换为Document格式"""
        if not os.path.exists(self.data_path):
            raise FileNotFoundError(f"找不到商品数据文件: {self.data_path}")

        try:
            # 读取Excel文件
            df = pd.read_excel(
                self.data_path,
                # 指定列名
                names=['goodsName', 'price', 'category', 'description'],
                # 如果第一行是标题，设置header=0；如果第一行就是数据，设置header=None
                header=0
            )

            documents = []
            # 遍历DataFrame的每一行
            for _, row in df.iterrows():
                # 将书籍信息组合成文本
                content = f"商品名：{row['goodsName']}\n价格：{row['price']}\n"
                content += f"类别：{row['category']}\n简介：{row['description']}\n"

                # 创建Document对象
                doc = Document(
                    page_content=content,
                    metadata={
                        "goodsName": row['goodsName'],
                        "price": row['price'],
                        "category": row['category']
                    }
                )
                documents.append(doc)

            return documents

        except Exception as e:
            raise Exception(f"读取Excel文件时发生错误: {str(e)}")