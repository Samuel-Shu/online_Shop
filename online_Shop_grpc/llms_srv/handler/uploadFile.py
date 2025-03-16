import proto.llms_pb2_grpc
import proto.llms_pb2
import grpc
import cmd


async def init_recommendation_system(file_path: str):
    """异步初始化推荐系统"""
    try:
        # 初始化图书加载器和推荐系统
        loader = cmd.GoodsLoader()
        documents = loader.load_Goods()

        global recommendation_system
        if recommendation_system is None:
            print("正在初始化推荐系统...")
            recommendation_system = cmd.GoodsRecommendationSystem()
            recommendation_system.initialize_vectorstore(documents)
        else:
            print("推荐系统已存在，正在重新初始化...")
            recommendation_system.initialize_vectorstore(documents)

        print("推荐系统初始化完成！")

    except Exception as e:
        print(f"初始化推荐系统时发生错误: {str(e)}")
class UploadFileServicer(llms_pb2_grpc.llmsServicer):
    def UploadFile(self, request, context):
        try:
            # 验证文件格式
            if not request.filename.endswith('.xlsx'):
                response = proto.llms_pb2.UploadFileResponse(
                    success=False,
                    message="文件格式错误，只支持excel文件，请重新上传"
                )
                return response

            # 保存上传的文件
            file_path = "./data/books.xlsx"
            os.makedirs("./data", exist_ok=True)

            # 如果文件已存在且内容相同，直接使用
            file_exists = os.path.exists(file_path)
            if not file_exists:
                # 保存新文件
                with open(file_path, "wb") as buffer:
                    shutil.copyfileobj(file.file, buffer)
                print("新数据文件已保存")
            else:
                print("检测到已存在的数据文件")

            # 异步启动初始化过程
            asyncio.create_task(init_recommendation_system(file_path))

            response = proto.llms_pb2.UploadFileResponse(
                success=True,
                message="文件上传成功，系统正在后台初始化"
            )
            return response

        except Exception as e:
            response = proto.llms_pb2.UploadFileResponse(
                success=False,
                message="文件上传失败"
            )
            return response
