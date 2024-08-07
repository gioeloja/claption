import torch
import torch.nn as nn
import torchvision.models as models
from torchvision.models.resnet import ResNet50_Weights
import numpy as np
import torchvision.transforms as transforms
from PIL import Image

device = torch.device("cpu")

"""
Convolutional Neural Network for image feature extraction.
Uses ResNet50 architecture as base and adds layers on top.

returns:
    features (dimensions [batch_size, height, width, channels])
"""
class Encoder(nn.Module):
    def __init__(self):         
        super(Encoder, self).__init__()
        self.device = 'cpu'

        self.resnet_base = models.resnet50(weights=ResNet50_Weights.DEFAULT)

        # Remove linear and pool layers (since we're not doing classification)
        modules = list(self.resnet_base.children())[:-2]
        self.resnet = nn.Sequential(*modules)

        for param in self.resnet.parameters():
            param.requires_grad = False
        
        # Fine-tuning only fine-tune convolutional blocks 2 through 4
        for c in list(self.resnet.children())[5:]:
            for p in c.parameters():
                p.requires_grad = True

        self.to(self.device)

    def forward(self, x):
        x = x.to(self.device)
        x = self.resnet(x)
        x = x.permute(0, 2, 3, 1) # changing dimensions from (batch_size, channels, height, width) to (batch_size, height, width, channels)
        return x
    
class Attention(nn.Module):
    """
    Attention Network.
    """

    def __init__(self, encoder_dim, decoder_dim, attention_dim):
        """
        Parameters:
            encoder_dim: feature size of encoded images
            decoder_dim: size of decoder's RNN
            attention_dim: size of the attention network
        """
        super(Attention, self).__init__()
        self.encoder_att = nn.Linear(encoder_dim, attention_dim)  # linear layer to transform encoded image
        self.decoder_att = nn.Linear(decoder_dim, attention_dim)  # linear layer to transform decoder's output
        self.full_att = nn.Linear(attention_dim, 1)  # linear layer to calculate values to be softmax-ed
        self.relu = nn.ReLU()
        self.softmax = nn.Softmax(dim=1)  # softmax layer to calculate weights

    def forward(self, encoder_out, decoder_hidden):
        """
        Parameters:
            encoder_out: encoded images, a tensor of dimension (batch_size, num_pixels, encoder_dim)
            decoder_hidden: previous decoder output, a tensor of dimension (batch_size, decoder_dim)

        returns: attention weighted encoding, weights
        """
        att1 = self.encoder_att(encoder_out)  # (batch_size, num_pixels, attention_dim)
        att2 = self.decoder_att(decoder_hidden)  # (batch_size, attention_dim)
        att = self.full_att(self.relu(att1 + att2.unsqueeze(1))).squeeze(2)  # (batch_size, num_pixels)
        alpha = self.softmax(att)  # (batch_size, num_pixels)
        attention_weighted_encoding = (encoder_out * alpha.unsqueeze(2)).sum(dim=1)  # (batch_size, encoder_dim)

        return attention_weighted_encoding, alpha


"""
Recurrent Neural Network for caption generation given features from the encoder
"""
class Decoder(nn.Module):
    def __init__(self, vocab_size, attention_dim, embedding_dim, decoder_dim, num_layers):
        super(Decoder, self).__init__()
        self.vocab_size = vocab_size
        self.attention_dim = attention_dim
        self.embedding_dim = embedding_dim
        self.decoder_dim = decoder_dim
        self.num_layers = num_layers

        self.attention = Attention(2048, decoder_dim, attention_dim)  # attention network
        
        # Embedding layer
        self.embedding = nn.Embedding(vocab_size, embedding_dim)
        
        # LSTM layers
        self.decode_step = nn.LSTMCell(embedding_dim + 2048, decoder_dim, bias=True)  # decoding LSTMCell
        self.lstm = nn.LSTM(embedding_dim, decoder_dim, num_layers, batch_first=True)

        self.init_h = nn.Linear(2048, decoder_dim)  # linear layer to find initial hidden state of LSTMCell (encoder_dim, decoder_dim)
        self.init_c = nn.Linear(2048, decoder_dim)  # linear layer to find initial cell state of LSTMCell (encoder_dim, decoder_dim)
        self.f_beta = nn.Linear(decoder_dim, 2048)  # linear layer to create a sigmoid-activated gate (encoder_dim, decoder_dim)
        self.sigmoid = nn.Sigmoid()

        # Output layer
        self.fc = nn.Linear(decoder_dim, vocab_size) # finds scores over vocabulary
    
    """
    Generate caption from encoder features. This function does not require
    the true caption for the image.

    Parameters:
        features: image features from encoder
        start_token: key for <start> in wordmap
        end_token: key for <end> in wordmap
    
    returns:
        array of wordmap tokens (array of nums)
    """
    def caption(self, features, max_length=10, start_token=114, end_token=115):
        start_token = torch.tensor([[start_token]], device=self.embedding.weight.device)

        # Flatten features
        encoder_dim = features.size(3)
        flattened_features = features.view(1, -1, encoder_dim)

        # Initialize the hidden and cell states
        mean_features = flattened_features.mean(dim=1)
        h = self.init_h(mean_features)
        c = self.init_c(mean_features)

        caption = [start_token.item()]
        alphas = []  # Store attention weights

        # Make caption
        for _ in range(max_length):
            embeddings = self.embedding(start_token).squeeze(1)
            attention_weight_encoding, alpha = self.attention(flattened_features, h)
            alpha = alpha.view(-1, features.size(1), features.size(1))
            gate = self.sigmoid(self.f_beta(h))
            attention_weight_encoding = gate * attention_weight_encoding
            h, c = self.decode_step(torch.cat([embeddings, attention_weight_encoding], dim=1), (h, c))
            preds = self.fc(h)

            # Get the predicted word
            pred_word_idx = preds.argmax(1)
            caption.append(pred_word_idx.item())
            alphas.append(alpha.cpu().detach().numpy().flatten())  # Flatten and save the attention weights

            # stop if end token is predicted
            if pred_word_idx == end_token:
                break

            start_token = torch.tensor([[pred_word_idx]], device=self.embedding.weight.device)

        # Convert alphas to a 2D NumPy array
        alphas = np.stack(alphas)
        
        return caption, alphas