/**********************************************************************
 * Copyright (c) 2015 Gregory Maxwell                                 *
 * Distributed under the MIT software license, see the accompanying   *
 * file COPYING or http://www.opensource.org/licenses/mit-license.php.*
 **********************************************************************/

#ifndef SECP256K1_MODULE_COMMITMENT_TESTS
#define SECP256K1_MODULE_COMMITMENT_TESTS

#include <string.h>

#include "group.h"
#include "scalar.h"
#include "testrand.h"
#include "util.h"

#include "include/secp256k1_commitment.h"

static void test_commitment_api(void) {
    secp256k1_pedersen_commitment commit;
    const secp256k1_pedersen_commitment *commit_ptr = &commit;
    unsigned char blind[32];
    unsigned char blind_out[32];
    const unsigned char *blind_ptr = blind;
    unsigned char *blind_out_ptr = blind_out;
    uint64_t val = secp256k1_rand32();

    secp256k1_context *none = secp256k1_context_create(SECP256K1_CONTEXT_NONE);
    secp256k1_context *sign = secp256k1_context_create(SECP256K1_CONTEXT_SIGN);
    secp256k1_context *vrfy = secp256k1_context_create(SECP256K1_CONTEXT_VERIFY);
    secp256k1_context *both = secp256k1_context_create(SECP256K1_CONTEXT_SIGN | SECP256K1_CONTEXT_VERIFY);
    int32_t ecount;

    secp256k1_context_set_error_callback(none, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_error_callback(sign, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_error_callback(vrfy, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_error_callback(both, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_illegal_callback(none, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_illegal_callback(sign, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_illegal_callback(vrfy, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_illegal_callback(both, counting_illegal_callback_fn, &ecount);

    secp256k1_rand256(blind);
    CHECK(secp256k1_pedersen_commit(none, &commit, blind, val, &secp256k1_generator_const_h, &secp256k1_generator_const_g) == 0);
    CHECK(ecount == 1);
    CHECK(secp256k1_pedersen_commit(vrfy, &commit, blind, val, &secp256k1_generator_const_h, &secp256k1_generator_const_g) == 0);
    CHECK(ecount == 2);
    CHECK(secp256k1_pedersen_commit(sign, &commit, blind, val, &secp256k1_generator_const_h, &secp256k1_generator_const_g) != 0);
    CHECK(ecount == 2);

    CHECK(secp256k1_pedersen_commit(sign, NULL, blind, val, &secp256k1_generator_const_h, &secp256k1_generator_const_g) == 0);
    CHECK(ecount == 3);
    CHECK(secp256k1_pedersen_commit(sign, &commit, NULL, val, &secp256k1_generator_const_h, &secp256k1_generator_const_g) == 0);
    CHECK(ecount == 4);
    CHECK(secp256k1_pedersen_commit(sign, &commit, blind, val, NULL, &secp256k1_generator_const_g) == 0);
    CHECK(ecount == 5);
    CHECK(secp256k1_pedersen_commit(sign, &commit, blind, val, &secp256k1_generator_const_h, NULL) == 0);
    CHECK(ecount == 6);

    CHECK(secp256k1_pedersen_blind_sum(none, blind_out, &blind_ptr, 1, 1) != 0);
    CHECK(ecount == 6);
    CHECK(secp256k1_pedersen_blind_sum(none, NULL, &blind_ptr, 1, 1) == 0);
    CHECK(ecount == 7);
    CHECK(secp256k1_pedersen_blind_sum(none, blind_out, NULL, 1, 1) == 0);
    CHECK(ecount == 8);
    CHECK(secp256k1_pedersen_blind_sum(none, blind_out, &blind_ptr, 0, 1) == 0);
    CHECK(ecount == 9);
    CHECK(secp256k1_pedersen_blind_sum(none, blind_out, &blind_ptr, 0, 0) != 0);
    CHECK(ecount == 9);

    CHECK(secp256k1_pedersen_commit(sign, &commit, blind, val, &secp256k1_generator_const_h, &secp256k1_generator_const_g) != 0);
    CHECK(secp256k1_pedersen_verify_tally(none, &commit_ptr, 1, &commit_ptr, 1) != 0);
    CHECK(secp256k1_pedersen_verify_tally(none, NULL, 0, &commit_ptr, 1) == 0);
    CHECK(secp256k1_pedersen_verify_tally(none, &commit_ptr, 1, NULL, 0) == 0);
    CHECK(secp256k1_pedersen_verify_tally(none, NULL, 0, NULL, 0) != 0);
    CHECK(ecount == 9);
    CHECK(secp256k1_pedersen_verify_tally(none, NULL, 1, &commit_ptr, 1) == 0);
    CHECK(ecount == 10);
    CHECK(secp256k1_pedersen_verify_tally(none, &commit_ptr, 1, NULL, 1) == 0);
    CHECK(ecount == 11);

    CHECK(secp256k1_pedersen_blind_generator_blind_sum(none, &val, &blind_ptr, &blind_out_ptr, 1, 0) != 0);
    CHECK(ecount == 11);
    CHECK(secp256k1_pedersen_blind_generator_blind_sum(none, &val, &blind_ptr, &blind_out_ptr, 1, 1) == 0);
    CHECK(ecount == 12);
    CHECK(secp256k1_pedersen_blind_generator_blind_sum(none, &val, &blind_ptr, &blind_out_ptr, 0, 0) == 0);
    CHECK(ecount == 13);
    CHECK(secp256k1_pedersen_blind_generator_blind_sum(none, NULL, &blind_ptr, &blind_out_ptr, 1, 0) == 0);
    CHECK(ecount == 14);
    CHECK(secp256k1_pedersen_blind_generator_blind_sum(none, &val, NULL, &blind_out_ptr, 1, 0) == 0);
    CHECK(ecount == 15);
    CHECK(secp256k1_pedersen_blind_generator_blind_sum(none, &val, &blind_ptr, NULL, 1, 0) == 0);
    CHECK(ecount == 16);

    secp256k1_context_destroy(none);
    secp256k1_context_destroy(sign);
    secp256k1_context_destroy(vrfy);
    secp256k1_context_destroy(both);
}

static void test_pedersen(void) {
    secp256k1_pedersen_commitment commits[19];
    const secp256k1_pedersen_commitment *cptr[19];
    unsigned char blinds[32*19];
    const unsigned char *bptr[19];
    secp256k1_scalar s;
    uint64_t values[19];
    int64_t totalv;
    int i;
    int inputs;
    int outputs;
    int total;
    inputs = (secp256k1_rand32() & 7) + 1;
    outputs = (secp256k1_rand32() & 7) + 2;
    total = inputs + outputs;
    for (i = 0; i < 19; i++) {
        cptr[i] = &commits[i];
        bptr[i] = &blinds[i * 32];
    }
    totalv = 0;
    for (i = 0; i < inputs; i++) {
        values[i] = secp256k1_rands64(0, INT64_MAX - totalv);
        totalv += values[i];
    }
    for (i = 0; i < outputs - 1; i++) {
        values[i + inputs] = secp256k1_rands64(0, totalv);
        totalv -= values[i + inputs];
    }
    values[total - 1] = totalv;

    for (i = 0; i < total - 1; i++) {
        random_scalar_order(&s);
        secp256k1_scalar_get_b32(&blinds[i * 32], &s);
    }
    CHECK(secp256k1_pedersen_blind_sum(ctx, &blinds[(total - 1) * 32], bptr, total - 1, inputs));
    for (i = 0; i < total; i++) {
        CHECK(secp256k1_pedersen_commit(ctx, &commits[i], &blinds[i * 32], values[i], &secp256k1_generator_const_h, &secp256k1_generator_const_g));
    }
    CHECK(secp256k1_pedersen_verify_tally(ctx, cptr, inputs, &cptr[inputs], outputs));
    CHECK(secp256k1_pedersen_verify_tally(ctx, &cptr[inputs], outputs, cptr, inputs));
    if (inputs > 0 && values[0] > 0) {
        CHECK(!secp256k1_pedersen_verify_tally(ctx, cptr, inputs - 1, &cptr[inputs], outputs));
    }
    random_scalar_order(&s);
    for (i = 0; i < 4; i++) {
        secp256k1_scalar_get_b32(&blinds[i * 32], &s);
    }
    values[0] = INT64_MAX;
    values[1] = 0;
    values[2] = 1;
    for (i = 0; i < 3; i++) {
        CHECK(secp256k1_pedersen_commit(ctx, &commits[i], &blinds[i * 32], values[i], &secp256k1_generator_const_h, &secp256k1_generator_const_g));
    }
    CHECK(secp256k1_pedersen_verify_tally(ctx, &cptr[0], 1, &cptr[0], 1));
    CHECK(secp256k1_pedersen_verify_tally(ctx, &cptr[1], 1, &cptr[1], 1));
}

#define MAX_N_GENS	30
void test_multiple_generators(void) {
    const size_t n_inputs = (secp256k1_rand32() % (MAX_N_GENS / 2)) + 1;
    const size_t n_outputs = (secp256k1_rand32() % (MAX_N_GENS / 2)) + 1;
    const size_t n_generators = n_inputs + n_outputs;
    unsigned char *generator_blind[MAX_N_GENS];
    unsigned char *pedersen_blind[MAX_N_GENS];
    secp256k1_generator generator[MAX_N_GENS];
    secp256k1_pedersen_commitment commit[MAX_N_GENS];
    const secp256k1_pedersen_commitment *commit_ptr[MAX_N_GENS];
    size_t i;
    int64_t total_value;
    uint64_t value[MAX_N_GENS];

    secp256k1_scalar s;

    unsigned char generator_seed[32];
    random_scalar_order(&s);
    secp256k1_scalar_get_b32(generator_seed, &s);
    /* Create all the needed generators */
    for (i = 0; i < n_generators; i++) {
        generator_blind[i] = (unsigned char*) malloc(32);
        pedersen_blind[i] = (unsigned char*) malloc(32);

        random_scalar_order(&s);
        secp256k1_scalar_get_b32(generator_blind[i], &s);
        random_scalar_order(&s);
        secp256k1_scalar_get_b32(pedersen_blind[i], &s);

        CHECK(secp256k1_generator_generate_blinded(ctx, &generator[i], generator_seed, generator_blind[i]));

        commit_ptr[i] = &commit[i];
    }

    /* Compute all the values -- can be positive or negative */
    total_value = 0;
    for (i = 0; i < n_outputs; i++) {
        value[n_inputs + i] = secp256k1_rands64(0, INT64_MAX - total_value);
        total_value += value[n_inputs + i];
    }
    for (i = 0; i < n_inputs - 1; i++) {
        value[i] = secp256k1_rands64(0, total_value);
        total_value -= value[i];
    }
    value[i] = total_value;

    /* Correct for blinding factors and do the commitments */
    CHECK(secp256k1_pedersen_blind_generator_blind_sum(ctx, value, (const unsigned char * const *) generator_blind, pedersen_blind, n_generators, n_inputs));
    for (i = 0; i < n_generators; i++) {
        CHECK(secp256k1_pedersen_commit(ctx, &commit[i], pedersen_blind[i], value[i], &generator[i], &secp256k1_generator_const_h));
    }

    /* Verify */
    CHECK(secp256k1_pedersen_verify_tally(ctx, &commit_ptr[0], n_inputs, &commit_ptr[n_inputs], n_outputs));

    /* Cleanup */
    for (i = 0; i < n_generators; i++) {
        free(generator_blind[i]);
        free(pedersen_blind[i]);
    }
}
#undef MAX_N_GENS

void run_commitment_tests(void) {
    int i;
    test_commitment_api();
    for (i = 0; i < 10*count; i++) {
        test_pedersen();
    }
    test_multiple_generators();
}

#endif
