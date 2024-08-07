import os
import torch
import json

def load_checkpoint(encoder, decoder, checkpoint_path):
    if os.path.exists(checkpoint_path):
        checkpoint = torch.load(checkpoint_path, map_location=torch.device('cpu'), weights_only=True)
        encoder.load_state_dict(checkpoint['encoder_state_dict'])
        decoder.load_state_dict(checkpoint['decoder_state_dict'])
        print("Checkpoint loaded")
    else:
        print("No checkpoint found.")
        

def load_json(file_path):
    res = {}

    with open(file_path, 'r') as file:
        res = json.load(file)

    return res