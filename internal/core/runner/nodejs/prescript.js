const argv = process.argv

// Check if seccomp should be completely disabled for Node.js
const disableSeccomp = process.env.DISABLE_NODEJS_SECCOMP === 'true';
// Check if seccomp should be skipped for debugging
const skipSeccomp = process.env.SKIP_SECCOMP === 'true';

if (disableSeccomp) {
    console.error('Node.js seccomp completely disabled via DISABLE_NODEJS_SECCOMP');
} else if (skipSeccomp) {
    console.error('Skipping seccomp initialization (debug mode)');
} else {
    try {
        // Try to require koffi and initialize seccomp
        const koffi = require('koffi')
        
        try {
            const lib = koffi.load('./var/sandbox/sandbox-nodejs/nodejs.so')
            const difySeccomp = lib.func('void DifySeccomp(int, int, bool)')
            
            const uid = parseInt(argv[2])
            const gid = parseInt(argv[3])
            const options = JSON.parse(argv[4])

            // Apply seccomp
            difySeccomp(uid, gid, options['enable_network'])
            
            console.error('Seccomp initialized successfully');
        } catch (libError) {
            console.error('Warning: Failed to load/initialize seccomp library:', libError.message);
            console.error('Continuing without seccomp protection...');
        }
    } catch (error) {
        console.error('Warning: Node.js sandbox initialization error:', error.message);
        console.error('Continuing without full sandbox protection...');
    }
}

