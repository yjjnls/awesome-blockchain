// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// secp256k1_context_create_sign_verify creates a context for signing and signature verification.
static secp256k1_context* secp256k1_context_create_sign_verify() {
	return secp256k1_context_create(SECP256K1_CONTEXT_SIGN | SECP256K1_CONTEXT_VERIFY);
}

// secp256k1_ecdsa_recover_pubkey recovers the public key of an encoded compact signature.
//
// Returns: 1: recovery was successful
//          0: recovery was not successful
// Args:    ctx:        pointer to a context object (cannot be NULL)
//  Out:    pubkey_out: the serialized 65-byte public key of the signer (cannot be NULL)
//  In:     sigdata:    pointer to a 65-byte signature with the recovery id at the end (cannot be NULL)
//          msgdata:    pointer to a 32-byte message (cannot be NULL)
static int secp256k1_ecdsa_recover_pubkey(
	const secp256k1_context* ctx,
	unsigned char *pubkey_out,
	const unsigned char *sigdata,
	const unsigned char *msgdata
) {
	secp256k1_ecdsa_recoverable_signature sig;
	secp256k1_pubkey pubkey;

	if (!secp256k1_ecdsa_recoverable_signature_parse_compact(ctx, &sig, sigdata, (int)sigdata[64])) {
		return 0;
	}
	if (!secp256k1_ecdsa_recover(ctx, &pubkey, &sig, msgdata)) {
		return 0;
	}
	size_t outputlen = 65;
	return secp256k1_ec_pubkey_serialize(ctx, pubkey_out, &outputlen, &pubkey, SECP256K1_EC_UNCOMPRESSED);
}

// secp256k1_pubkey_scalar_mul multiplies a point by a scalar in constant time.
//
// Returns: 1: multiplication was successful
//          0: scalar was invalid (zero or overflow)
// Args:    ctx:      pointer to a context object (cannot be NULL)
//  Out:    point:    the multiplied point (usually secret)
//  In:     point:    pointer to a 64-byte public point,
//                    encoded as two 256bit big-endian numbers.
//          scalar:   a 32-byte scalar with which to multiply the point
int secp256k1_pubkey_scalar_mul(const secp256k1_context* ctx, unsigned char *point, const unsigned char *scalar) {
	int ret = 0;
	int overflow = 0;
	secp256k1_fe feX, feY;
	secp256k1_gej res;
	secp256k1_ge ge;
	secp256k1_scalar s;
	ARG_CHECK(point != NULL);
	ARG_CHECK(scalar != NULL);
	(void)ctx;

	secp256k1_fe_set_b32(&feX, point);
	secp256k1_fe_set_b32(&feY, point+32);
	secp256k1_ge_set_xy(&ge, &feX, &feY);
	secp256k1_scalar_set_b32(&s, scalar, &overflow);
	if (overflow || secp256k1_scalar_is_zero(&s)) {
		ret = 0;
	} else {
		secp256k1_ecmult_const(&res, &ge, &s, 256);
		secp256k1_ge_set_gej(&ge, &res);
		/* Note: can't use secp256k1_pubkey_save here because it is not constant time. */
		secp256k1_fe_normalize(&ge.x);
		secp256k1_fe_normalize(&ge.y);
		secp256k1_fe_get_b32(point, &ge.x);
		secp256k1_fe_get_b32(point+32, &ge.y);
		ret = 1;
	}
	secp256k1_scalar_clear(&s);
	return ret;
}

void test_rangeproof() {
	size_t nbits, n_commits;
	//bench_bulletproof_rangeproof_t *data;
	/////////////////////////////////////////////////////////////
	bench_bulletproof_t odata;
	bench_bulletproof_rangeproof_t rp_data;

	odata.blind_gen = secp256k1_generator_const_g;
	odata.ctx = secp256k1_context_create(SECP256K1_CONTEXT_SIGN | SECP256K1_CONTEXT_VERIFY);
	odata.scratch = secp256k1_scratch_space_create(odata.ctx, 1024 * 1024 * 1024);
	odata.generators = secp256k1_bulletproof_generators_create(odata.ctx, &odata.blind_gen, 64 * 1024);

	rp_data.common = &odata;

	//run_rangeproof_test(&rp_data, 8, 1);
	/////////////////////////////////////////////////////////////
	nbits = 8;
	n_commits = 1;
	/////////////////////////////////////////////////////////////
	char str[64];

	(&rp_data)->nbits = nbits;
	(&rp_data)->n_commits = n_commits;
	(&rp_data)->common->iters = 100;

	(&rp_data)->common->n_proofs = 1;
	sprintf(str, "bulletproof_prove, %i, %i, 0, ", (int)nbits, (int) n_commits);
	
	//run_benchmark(str, bench_bulletproof_rangeproof_prove, bench_bulletproof_rangeproof_setup, bench_bulletproof_rangeproof_teardown, (void *)&rp_data, 5, 25);
	/////////////////////////////////////////////////////////////
        if (bench_bulletproof_rangeproof_setup != NULL) {
            bench_bulletproof_rangeproof_setup(&rp_data);
        }
        bench_bulletproof_rangeproof_prove(&rp_data);
        if (bench_bulletproof_rangeproof_teardown != NULL) {
            bench_bulletproof_rangeproof_teardown(&rp_data);
        }
	/////////////////////////////////////////////////////////////
        if (bench_bulletproof_rangeproof_setup != NULL) {
            bench_bulletproof_rangeproof_setup(&rp_data);
        }
        bench_bulletproof_rangeproof_verify(&rp_data);
        if (bench_bulletproof_rangeproof_teardown != NULL) {
            bench_bulletproof_rangeproof_teardown(&rp_data);
        }
	/////////////////////////////////////////////////////////////
}

static void counting_illegal_callback_fn(const char* str, void* data) {
    int32_t *p;
    (void)str;
    p = data;
    (*p)++;
}

typedef struct {
	size_t nbits;
	secp256k1_context *ctx_none; 
	secp256k1_context *ctx_both; 
	secp256k1_scratch *scratch;
	secp256k1_bulletproof_generators *gens;
	unsigned char *proof;
	size_t plen;
	uint64_t value;
	const unsigned char **blind_ptr;
	secp256k1_generator *value_gen;
	const unsigned char *blind;
	secp256k1_pedersen_commitment *pcommit;
} zkrp_t;

void myprint(char *message, zkrp_t *dt) {
    int i, len;
    len = (int) dt->plen; 
    printf("==================   %s   =================\n", message);
    printf("DT: %p\n", dt);
    printf("DT->nbits: %zd\n", dt->nbits);
    printf("DT->ctx_none: %p\n", dt->ctx_none);
    printf("DT->ctx_both: %p\n", dt->ctx_both);
    printf("DT->scratch: %p\n", dt->scratch);
    printf("DT->gens: %p\n", dt->gens);
    printf("DT->proof: %p\n", dt->proof);
    printf("DT->proof:\n");
    printf("[");
    for (i=0; i<len; i++) {
        printf("%d ", dt->proof[i]);
    }
    printf("]\n");
    printf("DT->plen: %zd\n", dt->plen);
    printf("DT->value: %lu\n", dt->value);
    printf("DT->blind: %p\n", dt->blind);
    printf("DT->pcommit: %p\n", dt->pcommit);
}


// setup should receive as input the range [A,B)
// and output a set of parameters
// - nbits means the interval is given by [0,2^nbits)
// - for now nbits must be in dt
void setup_rangeproof(zkrp_t *dt) {
    secp256k1_context *none = secp256k1_context_create(SECP256K1_CONTEXT_NONE);
    secp256k1_context *both = secp256k1_context_create(SECP256K1_CONTEXT_SIGN | SECP256K1_CONTEXT_VERIFY);
    secp256k1_scratch *scratch = secp256k1_scratch_space_create(both, 1024 * 1024);
    secp256k1_generator *value_gen = malloc(sizeof(secp256k1_generator));
    unsigned char blind[32] = "   i am not a blinding factor   ";

        dt->ctx_none = none;
        dt->ctx_both = both;
        dt->scratch = scratch;
        dt->value_gen = value_gen;

        dt->blind = (unsigned char*)malloc(sizeof(unsigned char) * 32);
        strncpy(blind, dt->blind, 32);

        dt->blind_ptr = malloc(sizeof(unsigned char*));
        dt->blind_ptr[0] = dt->blind;

    CHECK(secp256k1_generator_generate(dt->ctx_both, dt->value_gen, dt->blind) != 0);
}

void deallocate_memory(zkrp_t *dt) {
    // TODO realloc everything....
}


// prove should receive as input parameters and the commitment
// and output the proof
void prove_rangeproof(zkrp_t *dt) {
	printf("Prove result: %d\n", (secp256k1_bulletproof_rangeproof_prove(dt->ctx_both, dt->scratch, dt->gens, dt->proof, &(dt->plen), &(dt->value), NULL, dt->blind_ptr, 1, dt->value_gen, dt->nbits, dt->blind, NULL, 0) == 1));
}

// verify should receive as input parameters and proof
// and output true or false
int verify_rangeproof(zkrp_t *dt) {
	printf("Verification result: %d\n", (secp256k1_bulletproof_rangeproof_verify(dt->ctx_both, dt->scratch, dt->gens, dt->proof, dt->plen, NULL, dt->pcommit, 1, dt->nbits, dt->value_gen, NULL, 0) == 1));
	return 1;
}

// commit should receive parameters and value as input
// and output the commitment
void commit_rangeproof(zkrp_t *dt) {
    secp256k1_bulletproof_generators *gens;
    secp256k1_pedersen_commitment *pcommit = malloc(sizeof(secp256k1_pedersen_commitment));
    // TODO: value as input
    uint64_t value = 255;

    CHECK(secp256k1_pedersen_commit(dt->ctx_both, pcommit, dt->blind, value, dt->value_gen, &secp256k1_generator_const_h) != 0);
    
    gens = secp256k1_bulletproof_generators_create(dt->ctx_none, &secp256k1_generator_const_h, 256);
    CHECK(gens != NULL);
    
    dt->gens = gens;
    dt->proof = (unsigned char*)malloc(2000);
    dt->plen = 2000;
    dt->value = value;
    dt->pcommit = pcommit;
}

