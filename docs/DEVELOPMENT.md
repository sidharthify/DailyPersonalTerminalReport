# Development Guide

DPPR is designed to be easily extensible. This guide outlines how to contribute new modules or modify the core engine.

## Module Development

A module is simply a Python file in the `modules/` directory tree. It must have a `get_data(config)` function.

### Step-by-Step

1. **Create a file**: Place it in an existing category folder (e.g., `modules/custom/`) or create a new one.
2. **Implement `get_data`**:
   ```python
   def get_data(config):
       # Your data fetching logic
       return ["Result 1", "Result 2"]
   ```
3. **Handle Errors**: Use `try-except` blocks and return a useful error message as a string in the list if things fail.
4. **Avoid Hardcoding**: Use the `config` dictionary for any parameters.
5. **Add to Layout**: Test your module by adding it to `config.yaml`.

### Testing

Run the script with the `--no-print` flag to verify the output:
```bash
python main.py --no-print
```

Check the `daily_report.pdf` to see how your module's lines look.

---

## Engine Development

### pdf_generator.py
This file handles the rendering of text. If you want to change:
- **Font size or style**: Look in the `generate` function where `setFont` is called.
- **Section spacing**: Adjust the `y` coordinate decrementing logic.
- **Text cleaning**: The `_clean_text` method handles unicode normalization and character mapping.

### printer.py
This is a small wrapper around the `pycups` library. It handles finding the default printer and sending the PDF job.
