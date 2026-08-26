import torch

print("MPS tersedia:", torch.backends.mps.is_available())
print("MPS ter-build:", torch.backends.mps.is_built())
print("Versi PyTorch:", torch.__version__)