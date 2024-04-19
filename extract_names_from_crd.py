# Assisted by WCA for GP
# Latest GenAI contribution: granite-20B-code-instruct-v2 model
import os
import yaml

# Define the path to the directory containing the YAML files
path = 'config/manifests/bases'

# Create an empty list to store the names
names = []

# Loop through each file in the directory
for file in os.listdir(path):
    # Open the file
    with open(os.path.join(path, file), 'r') as f:
        # Read the contents of the file
        data = yaml.safe_load(f)

        # Extract the value of the .metadata.name field
        for crd in data['spec']['customresourcedefinitions']['owned']:
            name = crd['name']

            # Add the name to the names list
            names.append(name)

        # Add the name to the names list
        names.append(name)

# Sort the names alphabetically
names.sort()

# Print the names in the desired format
for name in names:
    print(f'    -\'*.{name}\'')
