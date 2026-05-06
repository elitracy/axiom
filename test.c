#include "./build/axiom.h"
#include "stdio.h"
#include <unistd.h>

int main() {
    char *result = AxiomInit("engine/cmd/axiom/initial_config.ax");
    printf("init: %s\n", result);

    while (1) {
        char *out = AxiomExecute("status");
        printf("status: %s\n", out);
        sleep(1);
    }

    return 0;
}
