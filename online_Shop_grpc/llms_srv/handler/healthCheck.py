import proto.llms_pb2_grpc
import proto.llms_pb2
import grpc


class HealthCheckServicer(llms_pb2_grpc.llmsServicer):
    def HealthCheck(self, request, context):
        response = proto.llms_pb2.HealthCheckResponse(
            success=recommendation_system is not None,
            message="llms初始化完成状态"
        )
        return response
