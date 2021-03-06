#!/usr/bin/env python
from grpc import aio
# from google.protobuf import empty_pb2
import asyncio
import logging
import sys
import signal
import grpc
import time
# pregenerated from proto file
import wedge_pb2_grpc
import wedge_pb2


def signal_handler(signal, frame):
    sys.exit(0)


class Downlink:
    def __init__(self, localhost):
        # 'localhost:30052'
        logging.info("Initialize Downlink object.")
        self.channel = grpc.insecure_channel(localhost)
        try:
            grpc.channel_ready_future(self.channel).result(timeout=1)
        except grpc.FutureTimeoutError:
            print('Error connecting to node: ', localhost)


class Wedge(wedge_pb2_grpc.WedgeServicer):

    def __init__(self, *args, **kwargs):
        pass

    def SetModel(self, request, context):
        # context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        # context.set_details('Method not implemented!')
        print(request)
        return wedge_pb2.Replay(ok=True)

    def SetDevice(self, request, context):
        # context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        # context.set_details('Method not implemented!')
        pass

    def SetValue(self, request, context):
        # context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        # context.set_details('Method not implemented!')
        pass

    def SetState(self, request, context):
        # context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        # context.set_details('Method not implemented!')
        return wedge_pb2.Replay(ok=True)

    def GetControl(self, request, context):
        print("Got request for Control, node: {}".format(request.node))
        time.sleep(10)
        print("Sending Control to  node: {}".format(request.node))
        update = wedge_pb2.UpdateState(
            state=wedge_pb2.State(
                id=1,
                data="99"
            ),
            device_id=1,
            value_id=1,
        )
        return wedge_pb2.Control(update=update)


async def downlinkLoop():
    # Once Node is connected it will send (SetModel) it's host and port.
    # Use that entry point to send back messages to Node/driver.
    pass


async def serve():
    server = aio.server()
    wedge_pb2_grpc.add_WedgeServicer_to_server(Wedge(), server)

    listen_addr = '[::]:50051'
    server.add_insecure_port(listen_addr)
    logging.info("Starting server on %s", listen_addr)
    await server.start()
    await server.wait_for_termination()


async def main():
    await asyncio.gather(downlinkLoop(), serve())

if __name__ == '__main__':
    signal.signal(signal.SIGINT, signal_handler)
    logging.basicConfig(level=logging.INFO)

    asyncio.run(main())
