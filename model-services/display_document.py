import matplotlib.pyplot as plt
from pdf2image import convert_from_path
from load_documents import PDFs

def display_pdf_images(images_list):
    num_images = len(images_list)
    num_rows = num_images // 5 + (1 if num_images % 5 > 0 else 0)
    fig, axes = plt.subplots(num_rows, 5, figsize=(20, 4 * num_rows))
    axes = axes.flatten()

    for i, img in enumerate(images_list):
        if i < len(axes):
            ax = axes[i]
            ax.imshow(img)
            ax.set_title(f"page {i+1}")
            ax.axis("off")

    for j in range(num_images, len(axes)):
        axes[j].axis('off')
    plt.tight_layout()
    plt.show()

for pdf in PDFs:
    pdf['images'] = convert_from_path(f"{pdf['title']}.pdf")
display_pdf_images(PDFs[0]["images"][:10])