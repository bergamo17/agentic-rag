import os
import sys
import json
import traceback
from document_builder import build_document

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

TEST_DIR = os.path.dirname(os.path.abspath(__file__))
OUTPUT_DIR = os.path.join(TEST_DIR, "_generated")
os.makedirs(OUTPUT_DIR, exist_ok=True)

def load_cases(filename):
    with open(os.path.join(TEST_DIR, filename)) as f:
        return json.load(f)

def run_cases(case: dict) -> tuple[bool, str]:
    expect = case.get("expect", "success")
    output_path = os.path.join(OUTPUT_DIR, f"{case['case_id']}.docx")
 
    try:
        build_document(
            theme_name=case["theme"],
            title=case["title"],
            sections=case["sections"],
            output_path=output_path,
        )
        if expect == "error":
            return False, "Expected an error, but build_document succeeded silently."
        if not os.path.exists(output_path):
            return False, "build_document returned success but no file was written."
        return True, f"OK -> {output_path}"

    except Exception as e:
        if expect == "success":
            return False, f"Unexpected exception: {type(e).__name__}: {e}"

        if isinstance(e, (KeyError, AttributeError, IndexError)):
            return False, (
                f"Raised {type(e).__name__} with a raw/unclear message: {e!r} "
                f"-- this should be a ValueError with a clear explanation instead."
            )

        return True, f"OK -> correctly raised {type(e).__name__}: {e}"

def run_suite(filename: str):
    print(f"\n=== {filename} ===")
    cases = load_cases(filename)
    passed_count = 0
 
    for case in cases:
        passed, message = run_cases(case)
        status = "PASS" if passed else "FAIL"
        print(f"[{status}] {case['case_id']}: {message}")
        if passed:
            passed_count += 1
 
    print(f"-- {passed_count}/{len(cases)} passed in {filename}")
    return passed_count, len(cases)

if __name__ == "__main__":
    total_passed = 0
    total_count = 0
 
    for suite_file in ["valid_themes.json", "sections_types.json", "edge_cases.json"]:
        p, c = run_suite(suite_file)
        total_passed += p
        total_count += c
 
    print(f"\n=== SUMMARY: {total_passed}/{total_count} passed ===")
    sys.exit(0 if total_passed == total_count else 1)