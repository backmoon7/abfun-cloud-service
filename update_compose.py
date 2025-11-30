import yaml

with open('docker-compose.yml', 'r') as f:
    data = yaml.safe_load(f)

video_service = data['services']['video-service']

# Add RabbitMQ dependency
if 'rabbitmq' not in video_service['depends_on']:
    video_service['depends_on']['rabbitmq'] = {'condition': 'service_started'}

# Add RabbitMQ URL env var
env = video_service['environment']
has_rabbit = False
for e in env:
    if e.startswith('RABBITMQ_URL='):
        has_rabbit = True
        break
if not has_rabbit:
    env.append('RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/')

with open('docker-compose.yml', 'w') as f:
    yaml.dump(data, f, default_flow_style=False, sort_keys=False)
