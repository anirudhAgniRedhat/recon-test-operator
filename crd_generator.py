import os
import yaml

# Define the output directory using the current working directory
output_dir = os.path.join(os.getcwd(), "generatedCRDS")
os.makedirs(output_dir, exist_ok=True)  # Create the directory if it doesn't exist

# Function to generate nested fields
def generate_nested_fields(level, max_level):
    if level > max_level:
        return {"type": "string"}  # Base type for leaf nodes

    # Recursively generate nested fields
    return {
        "type": "object",
        "properties": {
            f"nestedFieldLevel{level}": generate_nested_fields(level + 1, max_level)
        }
    }

# Function to generate a CRD definition
def generate_crd(index):
    crd = {
        "apiVersion": "apiextensions.k8s.io/v1",
        "kind": "CustomResourceDefinition",
        "metadata": {
            "name": f"recontests{index}.example.anirudh.io"
        },
        "spec": {
            "group": "example.anirudh.io",
            "versions": [{
                "name": "v1",
                "served": True,
                "storage": True,
                "schema": {
                    "openAPIV3Schema": {
                        "type": "object",
                        "properties": generate_nested_fields(1, 10)  # Nested up to 10 levels
                    }
                }
            }],
            "scope": "Namespaced",
            "names": {
                "plural": f"recontests{index}",
                "singular": f"recontest{index}",
                "kind": f"Recontest{str(index).capitalize()}",
                "shortNames": [f"rc{index}"]
            }
        }
    }
    return crd

# Generate and save CRDs
for i in range(1, 101):
    crd = generate_crd(i)
    file_name = os.path.join(output_dir, f"recontest{i}.yaml")
    with open(file_name, "w") as yaml_file:
        yaml.dump(crd, yaml_file, default_flow_style=False)

print(f"Generated 100 CRD YAML files in the {output_dir} directory.")
