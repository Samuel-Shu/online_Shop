from dotenv import load_dotenv
import os

# 加载环境变量
load_dotenv()

# OpenAI API配置
OPENAI_API_BASE = "https://ms-fc-7963d553-445c.api-inference.modelscope.cn/v1"
OPENAI_API_KEY = "077335bc-7b0e-42c3-8162-3d8dbce56f61"
MODEL_ID = "Qwen/Qwen2-7B-Instruct-GGUF"

# 向量数据库配置
CHROMA_PERSIST_DIRECTORY = "./data/chroma_db"

# embedding模型配置
EMBEDDING_MODEL_NAME = "paraphrase-multilingual-MiniLM-L12-v2"  # 使用默认多语言模型