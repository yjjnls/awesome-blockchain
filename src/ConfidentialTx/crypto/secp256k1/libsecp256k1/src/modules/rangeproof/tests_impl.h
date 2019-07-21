/**********************************************************************
 * Copyright (c) 2015 Gregory Maxwell                                 *
 * Distributed under the MIT software license, see the accompanying   *
 * file COPYING or http://www.opensource.org/licenses/mit-license.php.*
 **********************************************************************/

#ifndef SECP256K1_MODULE_RANGEPROOF_TESTS
#define SECP256K1_MODULE_RANGEPROOF_TESTS

#include <string.h>

#include "group.h"
#include "scalar.h"
#include "testrand.h"
#include "util.h"

#include "include/secp256k1_commitment.h"
#include "include/secp256k1_rangeproof.h"

static void test_rangeproof_api(const secp256k1_context *none, const secp256k1_context *sign, const secp256k1_context *vrfy, const secp256k1_context *both, const int32_t *ecount) {
    unsigned char proof[5134];
    unsigned char blind[32];
    secp256k1_pedersen_commitment commit;
    uint64_t vmin = secp256k1_rand32();
    uint64_t val = vmin + secp256k1_rand32();
    size_t len = sizeof(proof);
    /* we'll switch to dylan thomas for this one */
    const unsigned char message[68] = "My tears are like the quiet drift / Of petals from some magic rose;";
    size_t mlen = sizeof(message);
    const unsigned char ext_commit[72] = "And all my grief flows from the rift / Of unremembered skies and snows.";
    size_t ext_commit_len = sizeof(ext_commit);

    secp256k1_rand256(blind);
    CHECK(secp256k1_pedersen_commit(ctx, &commit, blind, val, &secp256k1_generator_const_h, &secp256k1_generator_const_g));

    CHECK(secp256k1_rangeproof_sign(none, proof, &len, vmin, &commit, blind, commit.data, 0, 0, val, message, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
    CHECK(*ecount == 1);
    CHECK(secp256k1_rangeproof_sign(sign, proof, &len, vmin, &commit, blind, commit.data, 0, 0, val, message, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
    CHECK(*ecount == 2);
    CHECK(secp256k1_rangeproof_sign(vrfy, proof, &len, vmin, &commit, blind, commit.data, 0, 0, val, message, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
    CHECK(*ecount == 3);
    CHECK(secp256k1_rangeproof_sign(both, proof, &len, vmin, &commit, blind, commit.data, 0, 0, val, message, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) != 0);
    CHECK(*ecount == 3);

    CHECK(secp256k1_rangeproof_sign(both, NULL, &len, vmin, &commit, blind, commit.data, 0, 0, val, message, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
    CHECK(*ecount == 4);
    CHECK(secp256k1_rangeproof_sign(both, proof, NULL, vmin, &commit, blind, commit.data, 0, 0, val, message, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
    CHECK(*ecount == 5);
    CHECK(secp256k1_rangeproof_sign(both, proof, &len, vmin, NULL, blind, commit.data, 0, 0, val, message, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
    CHECK(*ecount == 6);
    CHECK(secp256k1_rangeproof_sign(both, proof, &len, vmin, &commit, NULL, commit.data, 0, 0, val, message, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
    CHECK(*ecount == 7);
    CHECK(secp256k1_rangeproof_sign(both, proof, &len, vmin, &commit, blind, NULL, 0, 0, val, message, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
    CHECK(*ecount == 8);
    CHECK(secp256k1_rangeproof_sign(both, proof, &len, vmin, &commit, blind, commit.data, 0, 0, vmin - 1, message, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
    CHECK(*ecount == 8);
    CHECK(secp256k1_rangeproof_sign(both, proof, &len, vmin, &commit, blind, commit.data, 0, 0, val, NULL, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
    CHECK(*ecount == 9);
    CHECK(secp256k1_rangeproof_sign(both, proof, &len, vmin, &commit, blind, commit.data, 0, 0, val, NULL, 0, ext_commit, ext_commit_len, &secp256k1_generator_const_h) != 0);
    CHECK(*ecount == 9);
    CHECK(secp256k1_rangeproof_sign(both, proof, &len, vmin, &commit, blind, commit.data, 0, 0, val, NULL, 0, NULL, ext_commit_len, &secp256k1_generator_const_h) == 0);
    CHECK(*ecount == 10);
    CHECK(secp256k1_rangeproof_sign(both, proof, &len, vmin, &commit, blind, commit.data, 0, 0, val, NULL, 0, NULL, 0, &secp256k1_generator_const_h) != 0);
    CHECK(*ecount == 10);
    CHECK(secp256k1_rangeproof_sign(both, proof, &len, vmin, &commit, blind, commit.data, 0, 0, val, NULL, 0, NULL, 0, NULL) == 0);
    CHECK(*ecount == 11);

    CHECK(secp256k1_rangeproof_sign(both, proof, &len, vmin, &commit, blind, commit.data, 0, 0, val, message, mlen, ext_commit, ext_commit_len, &secp256k1_generator_const_h) != 0);
    {
        int exp;
        int mantissa;
        uint64_t min_value;
        uint64_t max_value;
        CHECK(secp256k1_rangeproof_info(none, &exp, &mantissa, &min_value, &max_value, proof, len) != 0);
        CHECK(exp == 0);
        CHECK(((uint64_t) 1 << mantissa) > val - vmin);
        CHECK(((uint64_t) 1 << (mantissa - 1)) <= val - vmin);
        CHECK(min_value == vmin);
        CHECK(max_value >= val);

        CHECK(secp256k1_rangeproof_info(none, NULL, &mantissa, &min_value, &max_value, proof, len) == 0);
        CHECK(*ecount == 12);
        CHECK(secp256k1_rangeproof_info(none, &exp, NULL, &min_value, &max_value, proof, len) == 0);
        CHECK(*ecount == 13);
        CHECK(secp256k1_rangeproof_info(none, &exp, &mantissa, NULL, &max_value, proof, len) == 0);
        CHECK(*ecount == 14);
        CHECK(secp256k1_rangeproof_info(none, &exp, &mantissa, &min_value, NULL, proof, len) == 0);
        CHECK(*ecount == 15);
        CHECK(secp256k1_rangeproof_info(none, &exp, &mantissa, &min_value, &max_value, NULL, len) == 0);
        CHECK(*ecount == 16);
        CHECK(secp256k1_rangeproof_info(none, &exp, &mantissa, &min_value, &max_value, proof, 0) == 0);
        CHECK(*ecount == 16);
    }
    {
        uint64_t min_value;
        uint64_t max_value;
        CHECK(secp256k1_rangeproof_verify(none, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 17);
        CHECK(secp256k1_rangeproof_verify(sign, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 18);
        CHECK(secp256k1_rangeproof_verify(vrfy, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) != 0);
        CHECK(*ecount == 18);

        CHECK(secp256k1_rangeproof_verify(vrfy, NULL, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 19);
        CHECK(secp256k1_rangeproof_verify(vrfy, &min_value, NULL, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 20);
        CHECK(secp256k1_rangeproof_verify(vrfy, &min_value, &max_value, NULL, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 21);
        CHECK(secp256k1_rangeproof_verify(vrfy, &min_value, &max_value, &commit, NULL, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 22);
        CHECK(secp256k1_rangeproof_verify(vrfy, &min_value, &max_value, &commit, proof, 0, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 22);
        CHECK(secp256k1_rangeproof_verify(vrfy, &min_value, &max_value, &commit, proof, len, NULL, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 23);
        CHECK(secp256k1_rangeproof_verify(vrfy, &min_value, &max_value, &commit, proof, len, NULL, 0, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 23);
        CHECK(secp256k1_rangeproof_verify(vrfy, &min_value, &max_value, &commit, proof, len, NULL, 0, NULL) == 0);
        CHECK(*ecount == 24);
    }
    {
        unsigned char blind_out[32];
        unsigned char message_out[68];
        uint64_t value_out;
        uint64_t min_value;
        uint64_t max_value;
        size_t message_len = sizeof(message_out);

        CHECK(secp256k1_rangeproof_rewind(none, blind_out, &value_out, message_out, &message_len, commit.data, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 25);
        CHECK(secp256k1_rangeproof_rewind(sign, blind_out, &value_out, message_out, &message_len, commit.data, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 26);
        CHECK(secp256k1_rangeproof_rewind(vrfy, blind_out, &value_out, message_out, &message_len, commit.data, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 27);
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, message_out, &message_len, commit.data, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) != 0);
        CHECK(*ecount == 27);

        CHECK(min_value == vmin);
        CHECK(max_value >= val);
        CHECK(value_out == val);
        CHECK(message_len == sizeof(message_out));
        CHECK(memcmp(message, message_out, sizeof(message_out)) == 0);

        CHECK(secp256k1_rangeproof_rewind(both, NULL, &value_out, message_out, &message_len, commit.data, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) != 0);
        CHECK(*ecount == 27);  /* blindout may be NULL */
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, NULL, message_out, &message_len, commit.data, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) != 0);
        CHECK(*ecount == 27);  /* valueout may be NULL */
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, NULL, &message_len, commit.data, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 28);
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, NULL, 0, commit.data, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) != 0);
        CHECK(*ecount == 28);
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, NULL, 0, NULL, &min_value, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 29);
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, NULL, 0, commit.data, NULL, &max_value, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 30);
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, NULL, 0, commit.data, &min_value, NULL, &commit, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 31);
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, NULL, 0, commit.data, &min_value, &max_value, NULL, proof, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 32);
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, NULL, 0, commit.data, &min_value, &max_value, &commit, NULL, len, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 33);
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, NULL, 0, commit.data, &min_value, &max_value, &commit, proof, 0, ext_commit, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 33);
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, NULL, 0, commit.data, &min_value, &max_value, &commit, proof, len, NULL, ext_commit_len, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 34);
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, NULL, 0, commit.data, &min_value, &max_value, &commit, proof, len, NULL, 0, &secp256k1_generator_const_h) == 0);
        CHECK(*ecount == 34);
        CHECK(secp256k1_rangeproof_rewind(both, blind_out, &value_out, NULL, 0, commit.data, &min_value, &max_value, &commit, proof, len, NULL, 0, NULL) == 0);
        CHECK(*ecount == 35);
    }
}

static void test_api(void) {
    secp256k1_context *none = secp256k1_context_create(SECP256K1_CONTEXT_NONE);
    secp256k1_context *sign = secp256k1_context_create(SECP256K1_CONTEXT_SIGN);
    secp256k1_context *vrfy = secp256k1_context_create(SECP256K1_CONTEXT_VERIFY);
    secp256k1_context *both = secp256k1_context_create(SECP256K1_CONTEXT_SIGN | SECP256K1_CONTEXT_VERIFY);
    int32_t ecount;
    int i;

    secp256k1_context_set_error_callback(none, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_error_callback(sign, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_error_callback(vrfy, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_error_callback(both, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_illegal_callback(none, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_illegal_callback(sign, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_illegal_callback(vrfy, counting_illegal_callback_fn, &ecount);
    secp256k1_context_set_illegal_callback(both, counting_illegal_callback_fn, &ecount);

    for (i = 0; i < count; i++) {
        ecount = 0;
        test_rangeproof_api(none, sign, vrfy, both, &ecount);
    }

    secp256k1_context_destroy(none);
    secp256k1_context_destroy(sign);
    secp256k1_context_destroy(vrfy);
    secp256k1_context_destroy(both);
}

static void test_borromean(void) {
    unsigned char e0[32];
    secp256k1_scalar s[64];
    secp256k1_gej pubs[64];
    secp256k1_scalar k[8];
    secp256k1_scalar sec[8];
    secp256k1_ge ge;
    secp256k1_scalar one;
    unsigned char m[32];
    size_t rsizes[8];
    size_t secidx[8];
    size_t nrings;
    size_t i;
    size_t j;
    int c;
    secp256k1_rand256_test(m);
    nrings = 1 + (secp256k1_rand32()&7);
    c = 0;
    secp256k1_scalar_set_int(&one, 1);
    if (secp256k1_rand32()&1) {
        secp256k1_scalar_negate(&one, &one);
    }
    for (i = 0; i < nrings; i++) {
        rsizes[i] = 1 + (secp256k1_rand32()&7);
        secidx[i] = secp256k1_rand32() % rsizes[i];
        random_scalar_order(&sec[i]);
        random_scalar_order(&k[i]);
        if(secp256k1_rand32()&7) {
            sec[i] = one;
        }
        if(secp256k1_rand32()&7) {
            k[i] = one;
        }
        for (j = 0; j < rsizes[i]; j++) {
            random_scalar_order(&s[c + j]);
            if(secp256k1_rand32()&7) {
                s[i] = one;
            }
            if (j == secidx[i]) {
                secp256k1_ecmult_gen(&ctx->ecmult_gen_ctx, &pubs[c + j], &sec[i]);
            } else {
                random_group_element_test(&ge);
                random_group_element_jacobian_test(&pubs[c + j],&ge);
            }
        }
        c += rsizes[i];
    }
    CHECK(secp256k1_borromean_sign(&ctx->ecmult_ctx, &ctx->ecmult_gen_ctx, e0, s, pubs, k, sec, rsizes, secidx, nrings, m, 32));
    CHECK(secp256k1_borromean_verify(&ctx->ecmult_ctx, NULL, e0, s, pubs, rsizes, nrings, m, 32));
    i = secp256k1_rand32() % c;
    secp256k1_scalar_negate(&s[i],&s[i]);
    CHECK(!secp256k1_borromean_verify(&ctx->ecmult_ctx, NULL, e0, s, pubs, rsizes, nrings, m, 32));
    secp256k1_scalar_negate(&s[i],&s[i]);
    secp256k1_scalar_set_int(&one, 1);
    for(j = 0; j < 4; j++) {
        i = secp256k1_rand32() % c;
        if (secp256k1_rand32() & 1) {
            secp256k1_gej_double_var(&pubs[i],&pubs[i], NULL);
        } else {
            secp256k1_scalar_add(&s[i],&s[i],&one);
        }
        CHECK(!secp256k1_borromean_verify(&ctx->ecmult_ctx, NULL, e0, s, pubs, rsizes, nrings, m, 32));
    }
}

static void test_rangeproof(void) {
    const uint64_t testvs[11] = {0, 1, 5, 11, 65535, 65537, INT32_MAX, UINT32_MAX, INT64_MAX - 1, INT64_MAX, UINT64_MAX};
    secp256k1_pedersen_commitment commit;
    secp256k1_pedersen_commitment commit2;
    unsigned char proof[5134 + 1]; /* One additional byte to test if trailing bytes are rejected */
    unsigned char blind[32];
    unsigned char blindout[32];
    unsigned char message[4096];
    size_t mlen;
    uint64_t v;
    uint64_t vout;
    uint64_t vmin;
    uint64_t minv;
    uint64_t maxv;
    size_t len;
    size_t i;
    size_t j;
    size_t k;
    /* Short message is a Simone de Beauvoir quote */
    const unsigned char message_short[120] = "When I see my own likeness in the depths of someone else's consciousness,  I always experience a moment of panic.";
    /* Long message is 0xA5 with a bunch of this quote in the middle */
    unsigned char message_long[3968];
    memset(message_long, 0xa5, sizeof(message_long));
    for (i = 1200; i < 3600; i += 120) {
        memcpy(&message_long[i], message_short, sizeof(message_short));
    }

    secp256k1_rand256(blind);
    for (i = 0; i < 11; i++) {
        v = testvs[i];
        CHECK(secp256k1_pedersen_commit(ctx, &commit, blind, v, &secp256k1_generator_const_h, &secp256k1_generator_const_g));
        for (vmin = 0; vmin < (i<9 && i > 0 ? 2 : 1); vmin++) {
            const unsigned char *input_message = NULL;
            size_t input_message_len = 0;
            /* vmin is always either 0 or 1; if it is 1, then we have no room for a message.
             * If it's 0, we use "minimum encoding" and only have room for a small message when
             * `testvs[i]` is >= 4; for a large message when it's >= 2^32. */
            if (vmin == 0 && i > 2) {
                input_message = message_short;
                input_message_len = sizeof(message_short);
            }
            if (vmin == 0 && i > 7) {
                input_message = message_long;
                input_message_len = sizeof(message_long);
            }
            len = 5134;
            CHECK(secp256k1_rangeproof_sign(ctx, proof, &len, vmin, &commit, blind, commit.data, 0, 0, v, input_message, input_message_len, NULL, 0, &secp256k1_generator_const_h));
            CHECK(len <= 5134);
            mlen = 4096;
            CHECK(secp256k1_rangeproof_rewind(ctx, blindout, &vout, message, &mlen, commit.data, &minv, &maxv, &commit, proof, len, NULL, 0, &secp256k1_generator_const_h));
            if (input_message != NULL) {
                CHECK(memcmp(message, input_message, input_message_len) == 0);
            }
            for (j = input_message_len; j < mlen; j++) {
                CHECK(message[j] == 0);
            }
            CHECK(mlen <= 4096);
            CHECK(memcmp(blindout, blind, 32) == 0);
            CHECK(vout == v);
            CHECK(minv <= v);
            CHECK(maxv >= v);
            len = 5134;
            CHECK(secp256k1_rangeproof_sign(ctx, proof, &len, v, &commit, blind, commit.data, -1, 64, v, NULL, 0, NULL, 0, &secp256k1_generator_const_h));
            CHECK(len <= 73);
            CHECK(secp256k1_rangeproof_rewind(ctx, blindout, &vout, NULL, NULL, commit.data, &minv, &maxv, &commit, proof, len, NULL, 0, &secp256k1_generator_const_h));
            CHECK(memcmp(blindout, blind, 32) == 0);
            CHECK(vout == v);
            CHECK(minv == v);
            CHECK(maxv == v);

            /* Check with a committed message */
            len = 5134;
            CHECK(secp256k1_rangeproof_sign(ctx, proof, &len, v, &commit, blind, commit.data, -1, 64, v, NULL, 0, message_short, sizeof(message_short), &secp256k1_generator_const_h));
            CHECK(len <= 73);
            CHECK(!secp256k1_rangeproof_rewind(ctx, blindout, &vout, NULL, NULL, commit.data, &minv, &maxv, &commit, proof, len, NULL, 0, &secp256k1_generator_const_h));
            CHECK(!secp256k1_rangeproof_rewind(ctx, blindout, &vout, NULL, NULL, commit.data, &minv, &maxv, &commit, proof, len, message_long, sizeof(message_long), &secp256k1_generator_const_h));
            CHECK(secp256k1_rangeproof_rewind(ctx, blindout, &vout, NULL, NULL, commit.data, &minv, &maxv, &commit, proof, len, message_short, sizeof(message_short), &secp256k1_generator_const_h));
            CHECK(memcmp(blindout, blind, 32) == 0);
            CHECK(vout == v);
            CHECK(minv == v);
            CHECK(maxv == v);
        }
    }
    secp256k1_rand256(blind);
    v = INT64_MAX - 1;
    CHECK(secp256k1_pedersen_commit(ctx, &commit, blind, v, &secp256k1_generator_const_h, &secp256k1_generator_const_g));
    for (i = 0; i < 19; i++) {
        len = 5134;
        CHECK(secp256k1_rangeproof_sign(ctx, proof, &len, 0, &commit, blind, commit.data, i, 0, v, NULL, 0, NULL, 0, &secp256k1_generator_const_h));
        CHECK(secp256k1_rangeproof_verify(ctx, &minv, &maxv, &commit, proof, len, NULL, 0, &secp256k1_generator_const_h));
        CHECK(len <= 5134);
        CHECK(minv <= v);
        CHECK(maxv >= v);
        /* Make sure it fails when validating with a committed message */
        CHECK(!secp256k1_rangeproof_verify(ctx, &minv, &maxv, &commit, proof, len, message_short, sizeof(message_short), &secp256k1_generator_const_h));
    }
    secp256k1_rand256(blind);
    {
        /*Malleability test.*/
        v = secp256k1_rands64(0, 255);
        CHECK(secp256k1_pedersen_commit(ctx, &commit, blind, v, &secp256k1_generator_const_h, &secp256k1_generator_const_g));
        len = 5134;
        CHECK(secp256k1_rangeproof_sign(ctx, proof, &len, 0, &commit, blind, commit.data, 0, 3, v, NULL, 0, NULL, 0, &secp256k1_generator_const_h));
        CHECK(len <= 5134);
        /* Test if trailing bytes are rejected. */
        proof[len] = v;
        CHECK(!secp256k1_rangeproof_verify(ctx, &minv, &maxv, &commit, proof, len + 1, NULL, 0, &secp256k1_generator_const_h));
        for (i = 0; i < len*8; i++) {
            proof[i >> 3] ^= 1 << (i & 7);
            CHECK(!secp256k1_rangeproof_verify(ctx, &minv, &maxv, &commit, proof, len, NULL, 0, &secp256k1_generator_const_h));
            proof[i >> 3] ^= 1 << (i & 7);
        }
        CHECK(secp256k1_rangeproof_verify(ctx, &minv, &maxv, &commit, proof, len, NULL, 0, &secp256k1_generator_const_h));
        CHECK(minv <= v);
        CHECK(maxv >= v);
    }
    memcpy(&commit2, &commit, sizeof(commit));
    for (i = 0; i < 10 * (size_t) count; i++) {
        int exp;
        int min_bits;
        v = secp256k1_rands64(0, UINT64_MAX >> (secp256k1_rand32()&63));
        vmin = 0;
        if ((v < INT64_MAX) && (secp256k1_rand32()&1)) {
            vmin = secp256k1_rands64(0, v);
        }
        secp256k1_rand256(blind);
        CHECK(secp256k1_pedersen_commit(ctx, &commit, blind, v, &secp256k1_generator_const_h, &secp256k1_generator_const_g));
        len = 5134;
        exp = (int)secp256k1_rands64(0,18)-(int)secp256k1_rands64(0,18);
        if (exp < 0) {
            exp = -exp;
        }
        min_bits = (int)secp256k1_rands64(0,64)-(int)secp256k1_rands64(0,64);
        if (min_bits < 0) {
            min_bits = -min_bits;
        }
        CHECK(secp256k1_rangeproof_sign(ctx, proof, &len, vmin, &commit, blind, commit.data, exp, min_bits, v, NULL, 0, NULL, 0, &secp256k1_generator_const_h));
        CHECK(len <= 5134);
        mlen = 4096;
        CHECK(secp256k1_rangeproof_rewind(ctx, blindout, &vout, message, &mlen, commit.data, &minv, &maxv, &commit, proof, len, NULL, 0, &secp256k1_generator_const_h));
        for (j = 0; j < mlen; j++) {
            CHECK(message[j] == 0);
        }
        CHECK(mlen <= 4096);
        CHECK(memcmp(blindout, blind, 32) == 0);
        CHECK(vout == v);
        CHECK(minv <= v);
        CHECK(maxv >= v);
        CHECK(secp256k1_rangeproof_rewind(ctx, blindout, &vout, NULL, NULL, commit.data, &minv, &maxv, &commit, proof, len, NULL, 0, &secp256k1_generator_const_h));
        memcpy(&commit2, &commit, sizeof(commit));
    }
    for (j = 0; j < 10; j++) {
        for (i = 0; i < 96; i++) {
            secp256k1_rand256(&proof[i * 32]);
        }
        for (k = 0; k < 128; k++) {
            len = k;
            CHECK(!secp256k1_rangeproof_verify(ctx, &minv, &maxv, &commit2, proof, len, NULL, 0, &secp256k1_generator_const_h));
        }
        len = secp256k1_rands64(0, 3072);
        CHECK(!secp256k1_rangeproof_verify(ctx, &minv, &maxv, &commit2, proof, len, NULL, 0, &secp256k1_generator_const_h));
    }
}

void run_rangeproof_tests(void) {
    int i;
    test_api();
    for (i = 0; i < 10*count; i++) {
        test_borromean();
    }
    test_rangeproof();
}

#endif
