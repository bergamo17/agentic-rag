import io
import os

import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt

from reportlab.lib import colors
from reportlab.lib.enums import TA_CENTER, TA_LEFT
from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import ParagraphStyle
from reportlab.lib.units import cm
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.platypus import (Image, KeepTogether, ListFlowable, ListItem,
                                Paragraph, SimpleDocTemplate, Spacer, Table,
                                TableStyle)

# ---------- Font ----------
FONT_DIR = os.environ.get("FONT_DIR", "/usr/share/fonts/truetype/dejavu/")
for name, f in [("DejaVu", "DejaVuSans.ttf"), ("DejaVu-Bold", "DejaVuSans-Bold.ttf"),
                ("DejaVu-Italic", "DejaVuSans-Oblique.ttf"), ("DejaVu-BoldItalic", "DejaVuSans-BoldOblique.ttf")]:
    pdfmetrics.registerFont(TTFont(name, FONT_DIR + f))
pdfmetrics.registerFontFamily("DejaVu", normal="DejaVu", bold="DejaVu-Bold",
                              italic="DejaVu-Italic", boldItalic="DejaVu-BoldItalic")

# ---------- Styles ----------
STYLES = {
    "title": ParagraphStyle("title", fontName="DejaVu-Bold", fontSize=18, leading=22,
                            alignment=TA_CENTER, spaceAfter=16),
    "h2": ParagraphStyle("h2", fontName="DejaVu-Bold", fontSize=13, leading=17,
                         spaceBefore=14, spaceAfter=6, keepWithNext=1),
    "body": ParagraphStyle("body", fontName="DejaVu", fontSize=10, leading=15,
                           alignment=TA_LEFT, spaceAfter=6),
}
STYLES["callout"] = ParagraphStyle(
    "callout", parent=STYLES["body"], backColor=colors.HexColor("#EEF4FF"),
    borderColor=colors.HexColor("#3B6FD4"), borderWidth=0.8,
    borderPadding=8, spaceBefore=8, spaceAfter=12)

# ---------- Story (otomatis terisi) ----------
_story = []

def _add(flowable):
    _story.append(flowable)
    return flowable

def title(text):   return _add(Paragraph(text, STYLES["title"]))
def h2(text):      return _add(Paragraph(text, STYLES["h2"]))
def body(text):    return _add(Paragraph(text, STYLES["body"]))
def spacer(h=8):   return _add(Spacer(1, h))
def callout(text): return _add(Paragraph(text, STYLES["callout"]))

def bullets(items):
    return _add(ListFlowable(
        [ListItem(Paragraph(t, STYLES["body"])) for t in items],
        bulletType="bullet", start="•", bulletFontName="DejaVu",
        bulletFontSize=10, leftIndent=18))

def table(rows, col_widths=None):
    t = Table([[Paragraph(str(c), STYLES["body"]) for c in r] for r in rows],
              colWidths=col_widths, repeatRows=1)
    t.setStyle(TableStyle([
        ("BACKGROUND", (0, 0), (-1, 0), colors.HexColor("#E5E7EB")),
        ("GRID", (0, 0), (-1, -1), 0.4, colors.grey),
        ("VALIGN", (0, 0), (-1, -1), "TOP")]))
    return _add(t)

def formula(latex, size=14):
    fig = plt.figure(figsize=(0.01, 0.01))
    fig.text(0, 0, f"${latex}$", fontsize=size)
    buf = io.BytesIO()
    fig.savefig(buf, format="png", dpi=300, bbox_inches="tight",
                pad_inches=0.08, transparent=True)
    plt.close(fig)
    buf.seek(0)
    img = Image(buf)
    img.drawWidth, img.drawHeight = img.imageWidth * 0.24, img.imageHeight * 0.24
    img.hAlign = "CENTER"
    return _add(img)

def figure(fig, width_cm=12):
    buf = io.BytesIO()
    fig.savefig(buf, format="png", dpi=200, bbox_inches="tight")
    plt.close(fig)
    buf.seek(0)
    img = Image(buf)
    ratio = img.imageHeight / img.imageWidth
    img.drawWidth, img.drawHeight = width_cm * cm, width_cm * cm * ratio
    img.hAlign = "CENTER"
    return _add(img)

# ---------- Build ----------
def build(path, title_text):
    try:
        SimpleDocTemplate(path, pagesize=A4, topMargin=2*cm, bottomMargin=2*cm,
                          leftMargin=2*cm, rightMargin=2*cm,
                          title=title_text).build(_story)
    finally:
        _story.clear()