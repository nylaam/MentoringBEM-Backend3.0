#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_VAULT_ITEMS 32
#define KEY_LEN 32
#define VAL_LEN 128

typedef struct {
    char key[KEY_LEN];
    char value[VAL_LEN];
    int active;
} VaultItem;

VaultItem vault[MAX_VAULT_ITEMS];

void init_service() {
    setvbuf(stdin, NULL, _IONBF, 0);
    setvbuf(stdout, NULL, _IONBF, 0);
    setvbuf(stderr, NULL, _IONBF, 0);
    memset(vault, 0, sizeof(vault));
}

void get_flag(char *buf, size_t max_len) {
    const char *flag_path = getenv("RSCTF_FLAG_FILE");
    if (!flag_path) {
        flag_path = "/flag";
    }
    FILE *f = fopen(flag_path, "r");
    if (f) {
        if (fgets(buf, max_len, f) == NULL) {
            strncpy(buf, "flag{dummy_flag_missing}", max_len);
        }
        buf[strcspn(buf, "\r\n")] = 0;
        fclose(f);
    } else {
        strncpy(buf, "flag{default_warmup_flag}", max_len);
    }
}

void handle_ping() {
    printf("PONG\n");
}

void handle_store(const char *key, const char *val) {
    for (int i = 0; i < MAX_VAULT_ITEMS; i++) {
        if (!vault[i].active || strcmp(vault[i].key, key) == 0) {
            strncpy(vault[i].key, key, KEY_LEN - 1);
            strncpy(vault[i].value, val, VAL_LEN - 1);
            vault[i].active = 1;
            printf("OK: Stored '%s'\n", key);
            return;
        }
    }
    printf("ERR: Vault full\n");
}

void handle_get(const char *key) {
    char flag_buf[VAL_LEN] = {0};
    if (strcmp(key, "flag") == 0 || strcmp(key, "secret") == 0) {
        get_flag(flag_buf, sizeof(flag_buf));
        printf("VALUE: %s\n", flag_buf);
        return;
    }
    for (int i = 0; i < MAX_VAULT_ITEMS; i++) {
        if (vault[i].active && strcmp(vault[i].key, key) == 0) {
            printf("VALUE: %s\n", vault[i].value);
            return;
        }
    }
    printf("ERR: Key not found\n");
}

int main() {
    init_service();
    printf("Welcome to Secure Vault\n");

    char line[256];
    while (fgets(line, sizeof(line), stdin)) {
        line[strcspn(line, "\r\n")] = 0;
        if (strlen(line) == 0) {
            continue;
        }

        if (strcmp(line, "PING") == 0) {
            handle_ping();
        } else if (strncmp(line, "STORE ", 6) == 0) {
            char key[KEY_LEN] = {0};
            char val[VAL_LEN] = {0};
            if (sscanf(line + 6, "%31s %127[^\n]", key, val) == 2) {
                handle_store(key, val);
            } else {
                printf("ERR: Usage: STORE <key> <value>\n");
            }
        } else if (strncmp(line, "GET ", 4) == 0) {
            char key[KEY_LEN] = {0};
            if (sscanf(line + 4, "%31s", key) == 1) {
                handle_get(key);
            } else {
                printf("ERR: Usage: GET <key>\n");
            }
        } else if (strcmp(line, "FLAG") == 0) {
            char flag_buf[VAL_LEN] = {0};
            get_flag(flag_buf, sizeof(flag_buf));
            printf("FLAG: %s\n", flag_buf);
        } else if (strcmp(line, "QUIT") == 0 || strcmp(line, "EXIT") == 0) {
            printf("BYE\n");
            break;
        } else {
            printf("ERR: Unknown command\n");
        }
    }

    return 0;
}
