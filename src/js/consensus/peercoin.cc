bool CheckStakeKernelHash(unsigned int nBits,
                          const CBlockHeader &blockFrom,
                          unsigned int nTxPrevOffset,
                          const CTransaction &txPrev,
                          const COutPoint &prevout,
                          unsigned int nTimeTx,
                          uint256 &hashProofOfStake,
                          bool fPrintProofOfStake)
{
    if (nTimeTx < txPrev.nTime)// Transaction timestamp violation
        return error("CheckStakeKernelHash() : nTime violation");

    unsigned int nTimeBlockFrom = blockFrom.GetBlockTime();
    if (nTimeBlockFrom + nStakeMinAge > nTimeTx)// Min age requirement
        return error("CheckStakeKernelHash() : min age violation");

    //目标值使用nBits
    CBigNum bnTargetPerCoinDay;
    bnTargetPerCoinDay.SetCompact(nBits);
    int64 nValueIn = txPrev.vout[prevout.n].nValue;
    // v0.3 protocol kernel hash weight starts from 0 at the 30-day min age
    // this change increases active coins participating the hash and helps
    // to secure the network when proof-of-stake difficulty is low
    int64 nTimeWeight = min((int64)nTimeTx - txPrev.nTime, (int64)STAKE_MAX_AGE) -
                        (IsProtocolV03(nTimeTx) ? nStakeMinAge : 0);
    //计算币龄，STAKE_MAX_AGE为90天
    CBigNum bnCoinDayWeight = CBigNum(nValueIn) * nTimeWeight / COIN / (24 * 60 * 60);
    // Calculate hash
    CDataStream ss(SER_GETHASH, 0);
    //权重修正因子
    uint64 nStakeModifier = 0;
    int nStakeModifierHeight = 0;
    int64 nStakeModifierTime = 0;
    if (IsProtocolV03(nTimeTx))// v0.3 protocol
    {
        if (!GetKernelStakeModifier(blockFrom.GetHash(), nTimeTx, nStakeModifier, nStakeModifierHeight, nStakeModifierTime, fPrintProofOfStake))
            return false;
        ss << nStakeModifier;
    }
    else// v0.2 protocol
    {
        ss << nBits;
    }

    //计算proofhash
    //即计算hash(nStakeModifier + txPrev.block.nTime + txPrev.offset + txPrev.nTime + txPrev.vout.n + nTime)
    ss << nTimeBlockFrom << nTxPrevOffset << txPrev.nTime << prevout.n << nTimeTx;
    hashProofOfStake = Hash(ss.begin(), ss.end());
    if (fPrintProofOfStake)
    {
        if (IsProtocolV03(nTimeTx))
            printf("CheckStakeKernelHash() : using modifier 0x%016" PRI64x " at height=%d timestamp=%s for block from height=%d timestamp=%s\n",
                   nStakeModifier, nStakeModifierHeight,
                   DateTimeStrFormat(nStakeModifierTime).c_str(),
                   mapBlockIndex[blockFrom.GetHash()]->nHeight,
                   DateTimeStrFormat(blockFrom.GetBlockTime()).c_str());
        printf("CheckStakeKernelHash() : check protocol=%s modifier=0x%016" PRI64x " nTimeBlockFrom=%u nTxPrevOffset=%u nTimeTxPrev=%u nPrevout=%u nTimeTx=%u hashProof=%s\n",
               IsProtocolV05(nTimeTx) ? "0.5" : (IsProtocolV03(nTimeTx) ? "0.3" : "0.2"),
               IsProtocolV03(nTimeTx) ? nStakeModifier : (uint64)nBits,
               nTimeBlockFrom, nTxPrevOffset, txPrev.nTime, prevout.n, nTimeTx,
               hashProofOfStake.ToString().c_str());
    }

    // Now check if proof-of-stake hash meets target protocol
    //判断是否满足proofhash < 币龄x目标值
    if (CBigNum(hashProofOfStake) > bnCoinDayWeight * bnTargetPerCoinDay)
        return false;
    if (fDebug && !fPrintProofOfStake)
    {
        if (IsProtocolV03(nTimeTx))
            printf("CheckStakeKernelHash() : using modifier 0x%016" PRI64x " at height=%d timestamp=%s for block from height=%d timestamp=%s\n",
                   nStakeModifier, nStakeModifierHeight,
                   DateTimeStrFormat(nStakeModifierTime).c_str(),
                   mapBlockIndex[blockFrom.GetHash()]->nHeight,
                   DateTimeStrFormat(blockFrom.GetBlockTime()).c_str());
        printf("CheckStakeKernelHash() : pass protocol=%s modifier=0x%016" PRI64x " nTimeBlockFrom=%u nTxPrevOffset=%u nTimeTxPrev=%u nPrevout=%u nTimeTx=%u hashProof=%s\n",
               IsProtocolV03(nTimeTx) ? "0.3" : "0.2",
               IsProtocolV03(nTimeTx) ? nStakeModifier : (uint64)nBits,
               nTimeBlockFrom, nTxPrevOffset, txPrev.nTime, prevout.n, nTimeTx,
               hashProofOfStake.ToString().c_str());
    }
    return true;
}

unsigned int static GetNextTargetRequired(const CBlockIndex *pindexLast, bool fProofOfStake)
{
    if (pindexLast == NULL)
        return bnProofOfWorkLimit.GetCompact();// genesis block

    const CBlockIndex *pindexPrev = GetLastBlockIndex(pindexLast, fProofOfStake);
    if (pindexPrev->pprev == NULL)
        return bnInitialHashTarget.GetCompact();// first block
    const CBlockIndex *pindexPrevPrev = GetLastBlockIndex(pindexPrev->pprev, fProofOfStake);
    if (pindexPrevPrev->pprev == NULL)
        return bnInitialHashTarget.GetCompact();// second block

    int64 nActualSpacing = pindexPrev->GetBlockTime() - pindexPrevPrev->GetBlockTime();

    // ppcoin: target change every block
    // ppcoin: retarget with exponential moving toward target spacing
    CBigNum bnNew;
    bnNew.SetCompact(pindexPrev->nBits);
    //STAKE_TARGET_SPACING为10分钟，即10 * 60
    //两个区块目标间隔时间即为10分钟
    int64 nTargetSpacing = fProofOfStake ? STAKE_TARGET_SPACING : min(nTargetSpacingWorkMax, (int64)STAKE_TARGET_SPACING * (1 + pindexLast->nHeight - pindexPrev->nHeight));
    //nTargetTimespan为1周，即7 * 24 * 60 * 60
    //nInterval为1008，即区块间隔为10分钟时，1周产生1008个区块
    int64 nInterval = nTargetTimespan / nTargetSpacing;
    //计算当前区块目标值
    bnNew *= ((nInterval - 1) * nTargetSpacing + nActualSpacing + nActualSpacing);
    bnNew /= ((nInterval + 1) * nTargetSpacing);

    if (bnNew > bnProofOfWorkLimit)
        bnNew = bnProofOfWorkLimit;

    return bnNew.GetCompact();
}

static bool CheckStakeKernelHashV2(CBlockIndex *pindexPrev,
                                   unsigned int nBits,
                                   unsigned int nTimeBlockFrom,
                                   const CTransaction &txPrev,
                                   const COutPoint &prevout,
                                   unsigned int nTimeTx,
                                   uint256 &hashProofOfStake,
                                   uint256 &targetProofOfStake,
                                   bool fPrintProofOfStake)
{
    if (nTimeTx < txPrev.nTime)// Transaction timestamp violation
        return error("CheckStakeKernelHash() : nTime violation");

    //目标值使用nBits
    CBigNum bnTarget;
    bnTarget.SetCompact(nBits);

    //计算币数x目标值
    int64_t nValueIn = txPrev.vout[prevout.n].nValue;
    CBigNum bnWeight = CBigNum(nValueIn);
    bnTarget *= bnWeight;

    targetProofOfStake = bnTarget.getuint256();

    //权重修正因子
    uint64_t nStakeModifier = pindexPrev->nStakeModifier;
    uint256 bnStakeModifierV2 = pindexPrev->bnStakeModifierV2;
    int nStakeModifierHeight = pindexPrev->nHeight;
    int64_t nStakeModifierTime = pindexPrev->nTime;

    //计算哈希值
    //即计算hash(nStakeModifier + txPrev.block.nTime + txPrev.nTime + txPrev.vout.hash + txPrev.vout.n + nTime)
    CDataStream ss(SER_GETHASH, 0);
    if (IsProtocolV3(nTimeTx))
        ss << bnStakeModifierV2;
    else
        ss << nStakeModifier << nTimeBlockFrom;
    ss << txPrev.nTime << prevout.hash << prevout.n << nTimeTx;
    hashProofOfStake = Hash(ss.begin(), ss.end());

    if (fPrintProofOfStake)
    {
        LogPrintf("CheckStakeKernelHash() : using modifier 0x%016x at height=%d timestamp=%s for block from timestamp=%s\n",
                  nStakeModifier, nStakeModifierHeight,
                  DateTimeStrFormat(nStakeModifierTime),
                  DateTimeStrFormat(nTimeBlockFrom));
        LogPrintf("CheckStakeKernelHash() : check modifier=0x%016x nTimeBlockFrom=%u nTimeTxPrev=%u nPrevout=%u nTimeTx=%u hashProof=%s\n",
                  nStakeModifier,
                  nTimeBlockFrom, txPrev.nTime, prevout.n, nTimeTx,
                  hashProofOfStake.ToString());
    }

    // Now check if proof-of-stake hash meets target protocol
    //判断是否满足proofhash < 币数x目标值
    if (CBigNum(hashProofOfStake) > bnTarget)
        return false;

    if (fDebug && !fPrintProofOfStake)
    {
        LogPrintf("CheckStakeKernelHash() : using modifier 0x%016x at height=%d timestamp=%s for block from timestamp=%s\n",
                  nStakeModifier, nStakeModifierHeight,
                  DateTimeStrFormat(nStakeModifierTime),
                  DateTimeStrFormat(nTimeBlockFrom));
        LogPrintf("CheckStakeKernelHash() : pass modifier=0x%016x nTimeBlockFrom=%u nTimeTxPrev=%u nPrevout=%u nTimeTx=%u hashProof=%s\n",
                  nStakeModifier,
                  nTimeBlockFrom, txPrev.nTime, prevout.n, nTimeTx,
                  hashProofOfStake.ToString());
    }

    return true;
}