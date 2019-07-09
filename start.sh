# First start consul
#consul agent -dev -bind 0.0.0.0 -client 0.0.0.0  &
# Now nomad agent
nomad agent -dev -config=/home/cneira/go/src/github.com/cneira/jail-task-driver/config.hcl -data-dir=/home/cneira/go/src/github.com/cneira/jail-task-driver -plugin-dir=/home/cneira/go/src/github.com/cneira/jail-task-driver/plugin -bind=0.0.0.0 
