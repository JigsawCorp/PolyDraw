import msgpack
import uuid
import time

sendMessage = None
uuidNul = "00000000-0000-0000-0000-000000000000"

def load(funcSendMessage):
    global sendMessage
    sendMessage = funcSendMessage

def execute():
    # You can add anything to reload
    #Example message
    """
    #Send message
    sendMessage(20, msgpack.packb(
        {"Message":"Hello ca va", "ChannelID":uuidNul}
    ))
    """
    """
    #Create channel
    sendMessage(26, msgpack.packb(
        {"ChannelName":"Boby"}
    ))
    """
    """
    #Destroy channel
    sendMessage(28, uuid.UUID("7040b302-6b1e-40f0-b9e8-12ee7301b2bf").bytes)
    """
    """
    #Join channel
    sendMessage(22, uuid.UUID("7040b302-6b1e-40f0-b9e8-12ee7301b2bf").bytes)
    """
    """
    #Quit channel
    sendMessage(24, uuid.UUID("7040b302-6b1e-40f0-b9e8-12ee7301b2bf").bytes)
    """
    """
    #Send message
    sendMessage(20, msgpack.packb(
        {"Message":"Hello ca va", "ChannelID":"7040b302-6b1e-40f0-b9e8-12ee7301b2bf"}
    ))
    """
    """
    #Send preview drawing
    sendMessage(36, uuid.UUID("7040b302-6b1e-40f0-b9e8-12ee7301b2bf").bytes)
    """
    """
    #Join group
    sendMessage(40, uuid.UUID("a4a97e0e-e9cf-4c9c-97e0-da5be7208933").bytes)
    """
    """
    #Leave group
    sendMessage(44, uuid.UUID("3688bc1b-6731-4a32-a515-d34d9cbef8e6").bytes)
    """
    """
    #Start match
    sendMessage(48, [])
    """
    """
    #Client is ready
    sendMessage(62, [])
    """