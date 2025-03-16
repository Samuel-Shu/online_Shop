import grpc
import proto.llms_pb2_grpc as gc
import handler.chat as ch
import handler.uploadFile as up
import handler.healthCheck as he


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    gc.add_llmsServicer_to_server(hd.ChatServicer, server)
    gc.add_llmsServicer_to_server(up.UploadFileServicer, server)
    gc.add_llmsServicer_to_server(he.HealthCheckServicer, server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started listening on port 50051")
    try:
        while True:
            time.sleep(86400)  # Keep the server running indefinitely
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
