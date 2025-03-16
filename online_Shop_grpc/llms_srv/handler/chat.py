import proto.llms_pb2_grpc
import proto.llms_pb2
import grpc


class ChatServicer(llms_pb2_grpc.llmsServicer):
    def SendMessage(self, request, context):
        """用户聊天接口"""
        try:
            if not recommendation_system:
                response = proto.llms_pb2.ChatMessageResponse(
                    success=False,
                    message="系统未初始化，请先上传商品数据"
                )
                return response

            user_input = request.get("ChatMessageRequest")
            if not user_input:
                response = proto.llms_pb2.ChatMessageResponse(
                    success=False,
                    message="消息不能为空"
                )
                return response

            # 创建自定义回调处理器
            class OutputCollector:
                def __init__(self):
                    self.output = []
                    self.raise_error = False
                    self.ignore_chat_model = False
                    self.ignore_llm = False

                def on_llm_new_token(self, token: str, **kwargs):
                    self.output.append(token)

                def on_llm_start(self, *args, **kwargs):
                    pass

                def on_llm_end(self, *args, **kwargs):
                    pass

                def on_llm_error(self, *args, **kwargs):
                    pass

                def on_chain_start(self, *args, **kwargs):
                    pass

                def on_chain_end(self, *args, **kwargs):
                    pass

                def on_chain_error(self, *args, **kwargs):
                    pass

                def on_chat_model_start(self, *args, **kwargs):
                    pass

            collector = OutputCollector()

            # 设置回调处理器
            original_callbacks = recommendation_system.llm.callbacks
            recommendation_system.llm.callbacks = [collector]

            try:
                # 判断是否是图书推荐查询
                is_recommendation = recommendation_system.is_goods_recommendation_query(user_input)

                # 根据查询类型处理消息
                if is_recommendation:
                    response = recommendation_system.get_recommendation(user_input)
                else:
                    response = recommendation_system.chat_prompt.pipe(recommendation_system.llm).invoke(
                        {"input": user_input}
                    )

                # 获取收集到的输出
                response_text = "".join(collector.output) if collector.output else response

                response = proto.llms_pb2.ChatMessageResponse(
                    success=True,
                    message=response_text
                )

                return response
            finally:
                # 恢复原始回调处理器
                recommendation_system.llm.callbacks = original_callbacks

        except Exception as e:

            response = proto.llms_pb2.ChatMessageResponse(
                success=False,
                message="聊天接口发生错误"
            )
            return response
