import json
import warnings

filename = 'result.txt'

with open(filename, 'r', encoding='utf-8') as f:
    results = f.read().strip().split('\n')

correct = 0
for result in results:
    
    result = json.loads(result)
    if result['answer'] == result['gold']:
        correct += 1

print('exact match:',correct / len(results))
